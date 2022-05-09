package common

import (
	"net/rpc"

	gop "github.com/hashicorp/go-plugin"

	"github.com/pipego/scheduler/plugin"
)

type Filter interface {
	Filter(*plugin.Args) Status
}

type Status struct {
	Error string
}

type FilterRPC struct {
	client *rpc.Client
}

func (n *FilterRPC) Filter(args *plugin.Args) Status {
	var resp Status
	if err := n.client.Call("Plugin.Filter", args, &resp); err != nil {
		panic(err)
	}
	return resp
}

type FilterRPCServer struct {
	Impl Filter
}

func (n *FilterRPCServer) Filter(args *plugin.Args, resp *Status) error {
	*resp = n.Impl.Filter(args)
	return nil
}

type FilterPlugin struct {
	Impl Filter
}

func (n *FilterPlugin) Server(*gop.MuxBroker) (interface{}, error) {
	return &FilterRPCServer{Impl: n.Impl}, nil
}

func (FilterPlugin) Client(b *gop.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &FilterRPC{client: c}, nil
}
