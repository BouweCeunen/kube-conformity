package main

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/stijndehaes/kube-conformity/kubeconformity"
)

var (
	master      string
	kubeconfig  string
	interval    time.Duration
	debug       bool
	version     string
	labels		[]string
	requests 	bool
	limits		bool
)

func init() {
	kingpin.Flag("master", "The address of the Kubernetes cluster to target").StringVar(&master)
	kingpin.Flag("kubeconfig", "Path to a kubeconfig file").StringVar(&kubeconfig)
	kingpin.Flag("interval", "Interval between conformity checks").Default("1h").DurationVar(&interval)
	kingpin.Flag("debug", "Enable debug logging.").BoolVar(&debug)
	kingpin.Flag("labels", "A list of labels that should be set on every pod in the cluster").StringsVar(&labels)
	kingpin.Flag("requestCheck", "Check if all pods have request filled in").BoolVar(&requests)
	kingpin.Flag("limitsCheck", "Check if all pods have limits filled in").BoolVar(&limits)
}

func main() {
	kingpin.Version(version)
	kingpin.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	client, err := newClient()
	if err != nil {
		log.Fatal(err)
	}

	rules := []kubeconformity.Rule{}

	if len(labels) != 0 {
		labelsRule := kubeconformity.LabelsFilledInRule{
			Labels: labels,
		}
		rules = append(rules, labelsRule)
	}

	if requests {
		requestRule := kubeconformity.RequestsFilledInRule{}
		rules = append(rules, requestRule)
	}

	if limits {
		limitsRule := kubeconformity.LimitsFilledInRule{}
		rules = append(rules, limitsRule)
	}

	kubeConformity := kubeconformity.New(
		client,
		log.StandardLogger(),
		rules,
	)

	for {
		err := kubeConformity.LogNonConformingPods()
		if err != nil {
			log.Fatal(err)
		}

		log.Debugf("Sleeping for %s...", interval)
		time.Sleep(interval)
	}
}

func newClient() (*kubernetes.Clientset, error) {
	if kubeconfig == "" {
		if _, err := os.Stat(clientcmd.RecommendedHomeFile); err == nil {
			kubeconfig = clientcmd.RecommendedHomeFile
		}
	}

	config, err := clientcmd.BuildConfigFromFlags(master, kubeconfig)
	if err != nil {
		return nil, err
	}

	log.Infof("Targeting cluster at %s", config.Host)

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
