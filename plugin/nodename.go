package main

import (
	"github.com/hashicorp/go-plugin"

	"github.com/pipego/plugin-filter/proto"
)

type NodeName struct{}

func (n *NodeName) Filter() string {
	return "TODO"
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
