package proto

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type NodeName interface {
	NodeName() string
}

type NodeNameRPC struct {
	client *rpc.Client
}

func (n *NodeNameRPC) NodeName() string {
	var resp string
	if err := n.client.Call("Plugin.NodeName", new(interface{}), &resp); err != nil {
		panic(err)
	}
	return resp
}

type NodeNameRPCServer struct {
	Impl NodeName
}

func (n *NodeNameRPCServer) NodeName(args interface{}, resp *string) error {
	*resp = n.Impl.NodeName()
	return nil
}

type NodeNamePlugin struct {
	Impl NodeName
}

func (n *NodeNamePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &NodeNameRPCServer{Impl: n.Impl}, nil
}

func (NodeNamePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &NodeNameRPC{client: c}, nil
}
