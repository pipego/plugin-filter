package main

import (
	"github.com/hashicorp/go-plugin"

	"github.com/pipego/plugin-filter/nodename/proto"
)

type NodeNameHello struct{}

func (n *NodeNameHello) NodeName() string {
	return "NodeName"
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "plugin-filter",
	MagicCookieValue: "nodename",
}

func main() {
	nodename := &NodeNameHello{}

	var pluginMap = map[string]plugin.Plugin{
		"nodename": &proto.NodeNamePlugin{Impl: nodename},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
