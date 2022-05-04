package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/pipego/plugin-filter/proto"
)

const (
	ErrReasonName = "NodeName: node(s) didn't match the requested node name"
)

type NodeName struct{}

func (n *NodeName) Filter(args *proto.Args) proto.Status {
	var status proto.Status

	if args.Task.NodeName != "" && args.Task.NodeName != args.Node.Name {
		status.Error = ErrReasonName
	}

	return status
}

// nolint:typecheck
func main() {
	config := plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "plugin-filter",
		MagicCookieValue: "plugin-filter",
	}

	pluginMap := map[string]plugin.Plugin{
		"NodeName": &proto.FilterPlugin{Impl: &NodeName{}},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: config,
		Plugins:         pluginMap,
	})
}
