# plugin-filter

[![Build Status](https://github.com/pipego/plugin-filter/workflows/ci/badge.svg?branch=main&event=push)](https://github.com/pipego/plugin-filter/actions?query=workflow%3Aci)
[![Go Report Card](https://goreportcard.com/badge/github.com/pipego/plugin-filter)](https://goreportcard.com/report/github.com/pipego/plugin-filter)
[![License](https://img.shields.io/github/license/pipego/plugin-filter.svg)](https://github.com/pipego/plugin-filter/blob/main/LICENSE)
[![Tag](https://img.shields.io/github/tag/pipego/plugin-filter.svg)](https://github.com/pipego/plugin-filter/tags)



## Introduction

*plugin-filter* is the filter plugin of [pipego](https://github.com/pipego) written in Go.



## Prerequisites

- Go >= 1.18.0



## Run

```bash
make lint
make build
./plugin-filter-test
```



## Docker



## Usage

- `plugin/nodeaffinity.go`: Implements node selectors and node affinity.
- `plugin/nodename.go`: Checks if a Task spec node name matches the current node.
- `plugin/noderesourcesfit.go`: Checks if the node has all the resources that the Task is requesting.
- `plugin/nodeunschedulable.go`: Filters out nodes that have .spec.unschedulable set to true.



## Settings



## License

Project License can be found [here](LICENSE).



## Reference

- [go-plugin](https://github.com/hashicorp/go-plugin)
- [kube-scheduler-config](https://kubernetes.io/docs/reference/scheduling/config)
- [kube-scheduler-interface](https://github.com/kubernetes/kubernetes/blob/master/pkg/scheduler/framework/interface.go)
- [kube-scheduler-plugins](https://github.com/kubernetes/kubernetes/blob/master/pkg/scheduler/framework/plugins)
