package proto

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type Filter interface {
	Filter(*Args) Status
}

type Args struct {
	Node Node
	Task Task
}

type Node struct {
	Name          string
	Unschedulable bool
}

type Task struct {
	NodeName               string
	ToleratesUnschedulable bool
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
