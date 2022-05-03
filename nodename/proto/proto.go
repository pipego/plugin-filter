package proto

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type NodeName interface {
	NodeName() string
}

type NodeNameRPC struct{ client *rpc.Client }

func (n *NodeNameRPC) NodeName() string {
	var resp string

	_ = n.client.Call("Plugin.NodeName", new(interface{}), &resp)

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

func (NodeNamePlugin) Client(_ *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &NodeNameRPC{client: c}, nil
}
