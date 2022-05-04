package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/pipego/plugin-filter/proto"
)

const (
	ErrReason = "node(s) didn't match the requested node name"
)

type NodeName struct{}

func (n *NodeName) Filter(args *proto.Args) proto.Status {
	var status proto.Status

	if args.Task.NodeName != "" && args.Task.NodeName != args.Node.Name {
		status.Error = ErrReason
	}

	return status
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "plugin-filter",
	MagicCookieValue: "plugin-filter",
}

func main() {
	n := &NodeName{}

	var pluginMap = map[string]plugin.Plugin{
		"NodeName": &proto.FilterPlugin{Impl: n},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
