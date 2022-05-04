package proto

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

const (
	// ResourceCPU CPU, in cores. (500m = .5 cores)
	ResourceCPU = "cpu"
	// ResourceMemory Memory, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourceMemory = "memory"
	// ResourceStorage Volume size, in bytes (e,g. 5Gi = 5GiB = 5 * 1024 * 1024 * 1024)
	ResourceStorage = "storage"
)

type Args struct {
	Node Node
	Task Task
}

type Node struct {
	AllocatableResource Resource `json:"allocatableResource"`
	Label               Label    `json:"label"`
	Name                string   `json:"name"`
	RequestedResource   Resource `json:"requestedResource"`
	Unschedulable       bool     `json:"unschedulable"`
}

type Task struct {
	NodeName               string   `json:"nodeName"`
	NodeSelector           Selector `json:"nodeSelector"`
	RequestedResource      Resource `json:"requestedResource"`
	ToleratesUnschedulable bool     `json:"toleratesUnschedulable"`
}

type Label map[string]string

type Resource struct {
	MilliCPU int64 `json:"milliCPU"`
	Memory   int64 `json:"memory"`
	Storage  int64 `json:"storage"`
}

type Selector map[string][]string

type Filter interface {
	Filter(*Args) Status
}

type Status struct {
	Error string
}

type FilterRPC struct {
	client *rpc.Client
}

func (n *FilterRPC) Filter(args *Args) Status {
	var resp Status
	if err := n.client.Call("Plugin.Filter", args, &resp); err != nil {
		panic(err)
	}
	return resp
}

type FilterRPCServer struct {
	Impl Filter
}

func (n *FilterRPCServer) Filter(args *Args, resp *Status) error {
	*resp = n.Impl.Filter(args)
	return nil
}

type FilterPlugin struct {
	Impl Filter
}

func (n *FilterPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &FilterRPCServer{Impl: n.Impl}, nil
}

func (FilterPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &FilterRPC{client: c}, nil
}
