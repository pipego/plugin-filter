package main

import (
	"github.com/hashicorp/go-plugin"

	"github.com/pipego/plugin-filter/nodename/proto"
)

type NodeName struct{}

func (n *NodeName) NodeName() string {
	return "TODO"
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "plugin-filter",
	MagicCookieValue: "NodeName",
}

func main() {
	n := &NodeName{}

	var pluginMap = map[string]plugin.Plugin{
		"NodeName": &proto.NodeNamePlugin{Impl: n},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
