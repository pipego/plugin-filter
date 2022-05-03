package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/pipego/plugin-filter/test/proto"
)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "test",
		Output: os.Stderr,
		Level:  hclog.Error,
	})

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command("./bin/plugin-filter-nodename"),
		Logger:          logger,
	})
	defer client.Kill()

	rpcClient, _ := client.Client()

	raw, _ := rpcClient.Dispense("nodename")
	nodename := raw.(proto.NodeName)

	fmt.Println(nodename.NodeName())
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

var pluginMap = map[string]plugin.Plugin{
	"nodename": &proto.NodeNamePlugin{},
}
