# Kube conformity

[![Build Status](https://travis-ci.org/stijndehaes/kube-conformity.svg?branch=master)](https://travis-ci.org/stijndehaes/kube-conformity)
[![Coverage Status](https://coveralls.io/repos/github/stijndehaes/kube-conformity/badge.svg?branch=master)](https://coveralls.io/github/stijndehaes/kube-conformity?branch=master)
[![Docker Hub](https://img.shields.io/docker/build/sdehaes/kube-conformity.svg)](https://hub.docker.com/r/sdehaes/kube-conformity/)
[![GitHub release](https://img.shields.io/github/release/stijndehaes/kube-conformity.svg)](https://github.com/stijndehaes/kube-conformity/releases)

This project looks for pods that in your kubernetes cluster that are breaking conformity rules.

At this moment there are three rules defined:

* Labels filled in: Takes a list of labels and check if pods have these labels defined
* Resource requests filled in: Checks all pods if they have resource requests filled in
* Limits requests filled in: Checks all pods if they have limits requests filled in

The rules are configured using a yaml config. An example of this config is:

```yaml
interval: 1h
labels_filled_in_rules:
- name: Check if label app is active on every pod
  labels:
  - app
  filter:
    namespaces: "!kube-system"
limits_filled_in_rules:
- name: Checks if limits are filled in everywhere
  filter:
    namespaces: "!kube-system"
requests_filled_in_rules:
- name: Checks if requests are filled in everywhere
  filter:
    namespaces: "!kube-system"
```

To run this project you can use the provided docker image.
An example deployment can be found in the examples folder.
If you have an RBAC enabled cluster the minimum clusterrole is defined in the examples folder.