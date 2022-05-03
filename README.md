# plugin-filter

[![License](https://img.shields.io/github/license/pipego/plugin-filter.svg)](https://github.com/pipego/plugin-filter/blob/main/LICENSE)



## Introduction

*plugin-filter* is the filter plugin of [pipego](https://github.com/pipego) written in Go.



## Prerequisites

- Go >= 1.17.0



## Run

```bash
# Template
go env -w GOPROXY=https://goproxy.cn,direct
CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w" -o plugin/filter-nodename plugin/nodename.go
```



```bash
# Test
go env -w GOPROXY=https://goproxy.cn,direct
CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w" -o plugin-filter-test main.go
./plugin-filter-test
```



## Docker



## Usage

- `NodeName`: Checks if a Task spec node name matches the current node.
- `NodeResourcesFit`: Checks if the node has all the resources that the Task is requesting.
- `NodeSelector`: Checks if a Task spec node label matches the current node.
- `NodeUnschedulable`: Filters out nodes that have .spec.unschedulable set to true.



## Settings



## License

Project License can be found [here](LICENSE).



## Reference

- [go-plugin](https://github.com/hashicorp/go-plugin)
- [kube-scheduler-config](https://kubernetes.io/docs/reference/scheduling/config)
- [kube-scheduler-interface](https://github.com/kubernetes/kubernetes/blob/master/pkg/scheduler/framework/interface.go)
- [kube-scheduler-plugins](https://github.com/kubernetes/kubernetes/blob/master/pkg/scheduler/framework/plugins)
