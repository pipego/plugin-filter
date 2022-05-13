package common

import (
	"net/rpc"

	gop "github.com/hashicorp/go-plugin"

	"github.com/pipego/scheduler/plugin"
)

type FilterRPC struct {
	client *rpc.Client
}

func (n *FilterRPC) Run(args *plugin.Args) plugin.FilterResult {
	var resp plugin.FilterResult
	if err := n.client.Call("Plugin.Run", args, &resp); err != nil {
		panic(err)
	}
	return resp
}

type FilterRPCServer struct {
	Impl plugin.FilterPlugin
}

func (n *FilterRPCServer) Run(args *plugin.Args, resp *plugin.FilterResult) error {
	*resp = n.Impl.Run(args)
	return nil
}

type FilterPlugin struct {
	Impl plugin.FilterPlugin
}

func (n *FilterPlugin) Server(*gop.MuxBroker) (interface{}, error) {
	return &FilterRPCServer{Impl: n.Impl}, nil
}

func (FilterPlugin) Client(b *gop.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &FilterRPC{client: c}, nil
}
