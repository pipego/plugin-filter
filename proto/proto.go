package proto

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type Filter interface {
	Filter() string
}

type FilterRPC struct {
	client *rpc.Client
}

func (n *FilterRPC) Filter() string {
	var resp string
	if err := n.client.Call("Plugin.Filter", new(interface{}), &resp); err != nil {
		panic(err)
	}
	return resp
}

type FilterRPCServer struct {
	Impl Filter
}

func (n *FilterRPCServer) Filter(args interface{}, resp *string) error {
	*resp = n.Impl.Filter()
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
