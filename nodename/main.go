package main

import (
	"os"
	"plugin"

	"github.com/pipego/plugin-filter/nodename/proto"
)

type NodeNameHello struct {
	logger hclog.Logger
}

func (n *NodeNameHello) NodeName() string {
	n.logger.Debug("message from NodeNameHello.NodeName")
	return "Hello"
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	nodename := &NodeNameHello{
		logger: logger,
	}

	var pluginMap = map[string]plugin.Plugin{
		"nodename": &proto.NodeNamePlugin{Impl: nodename},
	}

	logger.Debug("message from plugin", "foo", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
