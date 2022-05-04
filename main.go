package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/pipego/plugin-filter/proto"
)

type config struct {
	args *proto.Args
	name string
	path string
}

var (
	configs = []config{
		// Plugin: NodeName
		{
			args: &proto.Args{
				Node: proto.Node{
					Name: "Node",
				},
				Task: proto.Task{
					NodeName: "Task",
				},
			},
			name: "NodeName",
			path: "./plugin/filter-nodename",
		},
		// Plugin: NodeResourcesfit
		{
			args: &proto.Args{},
			name: "NodeResourcesFit",
			path: "./plugin/filter-noderesourcesfit",
		},
		// Plugin: NodeSelector
		{
			args: &proto.Args{},
			name: "NodeSelector",
			path: "./plugin/filter-nodeselector",
		},
		// Plugin: NodeUnschedulable
		{
			args: &proto.Args{
				Node: proto.Node{
					Unschedulable: true,
				},
				Task: proto.Task{
					ToleratesUnschedulable: false,
				},
			},
			name: "NodeUnschedulable",
			path: "./plugin/filter-nodeunschedulable",
		},
	}
)

func main() {
	for _, item := range configs {
		status, _ := helper(item.path, item.name, item.args)
		fmt.Println(status.Error)
	}
}

func helper(path, name string, args *proto.Args) (proto.Status, error) {
	config := plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "plugin-filter",
		MagicCookieValue: "plugin-filter",
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin-filter",
		Output: os.Stderr,
		Level:  hclog.Error,
	})

	plugins := map[string]plugin.Plugin{
		name: &proto.FilterPlugin{},
	}

	client := plugin.NewClient(&plugin.ClientConfig{
		Cmd:             exec.Command(path),
		HandshakeConfig: config,
		Logger:          logger,
		Plugins:         plugins,
	})
	defer client.Kill()

	rpcClient, _ := client.Client()
	raw, _ := rpcClient.Dispense(name)
	n := raw.(proto.Filter)
	status := n.Filter(args)

	return status, nil
}
