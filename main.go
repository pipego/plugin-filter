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
	// Plugin: NodeName
	args := &proto.Args{
		Node: proto.Node{
			Name: "Node",
		},
		Task: proto.Task{
			NodeName: "Task",
		},
	}

	status, _ := helper("./plugin/filter-nodename", "NodeName", args)
	fmt.Println(status.Error)

	// Plugin: NodeUnschedulable
	args = &proto.Args{
		Node: proto.Node{
			Unschedulable: true,
		},
		Task: proto.Task{
			ToleratesUnschedulable: false,
		},
	}

	status, _ = helper("./plugin/filter-nodeunschedulable", "NodeUnschedulable", args)
	fmt.Println(status.Error)
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
