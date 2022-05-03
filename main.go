package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/pipego/plugin-filter/proto"
)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin-filter",
		Output: os.Stderr,
		Level:  hclog.Error,
	})

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command("./plugin/filter-nodename"),
		Logger:          logger,
	})
	defer client.Kill()

	rpcClient, _ := client.Client()

	raw, _ := rpcClient.Dispense("NodeName")
	n := raw.(proto.Filter)

	fmt.Println(n.Filter())
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "plugin-filter",
	MagicCookieValue: "plugin-filter",
}

var pluginMap = map[string]plugin.Plugin{
	"NodeName": &proto.FilterPlugin{},
}
