package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	gop "github.com/hashicorp/go-plugin"

	"github.com/pipego/plugin-filter/common"
	"github.com/pipego/scheduler/plugin"
)

type config struct {
	args *plugin.Args
	name string
	path string
}

var (
	configs = []config{
		// Plugin: NodeName
		{
			args: &plugin.Args{
				Node: plugin.Node{
					Name: "Node",
				},
				Task: plugin.Task{
					NodeName: "Node",
				},
			},
			name: "NodeName",
			path: "./plugin/filter-nodename",
		},
		{
			args: &plugin.Args{
				Node: plugin.Node{
					Name: "Node",
				},
				Task: plugin.Task{
					NodeName: "Task",
				},
			},
			name: "NodeName",
			path: "./plugin/filter-nodename",
		},
		// Plugin: NodeResourcesFit
		{
			args: &plugin.Args{
				Node: plugin.Node{
					AllocatableResource: plugin.Resource{
						MilliCPU: 400,
					},
					RequestedResource: plugin.Resource{
						MilliCPU: 200,
					},
				},
				Task: plugin.Task{
					RequestedResource: plugin.Resource{
						MilliCPU: 100,
					},
				},
			},
			name: "NodeResourcesFit",
			path: "./plugin/filter-noderesourcesfit",
		},
		{
			args: &plugin.Args{
				Node: plugin.Node{
					AllocatableResource: plugin.Resource{
						MilliCPU: 400,
					},
					RequestedResource: plugin.Resource{
						MilliCPU: 200,
					},
				},
				Task: plugin.Task{
					RequestedResource: plugin.Resource{
						MilliCPU: 500,
					},
				},
			},
			name: "NodeResourcesFit",
			path: "./plugin/filter-noderesourcesfit",
		},
		// Plugin: NodeAffinity
		{
			args: &plugin.Args{
				Node: plugin.Node{
					Label: plugin.Label{
						"disktype": "ssd",
					},
				},
				Task: plugin.Task{
					NodeSelector: plugin.Selector{
						"disktype": []string{"ssd"},
					},
				},
			},
			name: "NodeAffinity",
			path: "./plugin/filter-nodeaffinity",
		},
		{
			args: &plugin.Args{
				Node: plugin.Node{
					Label: plugin.Label{
						"disktype": "ssd",
					},
				},
				Task: plugin.Task{
					NodeSelector: plugin.Selector{
						"disktype": []string{"hdd"},
					},
				},
			},
			name: "NodeAffinity",
			path: "./plugin/filter-nodeaffinity",
		},
		// Plugin: NodeUnschedulable
		{
			args: &plugin.Args{
				Node: plugin.Node{
					Unschedulable: true,
				},
				Task: plugin.Task{
					ToleratesUnschedulable: true,
				},
			},
			name: "NodeUnschedulable",
			path: "./plugin/filter-nodeunschedulable",
		},
		{
			args: &plugin.Args{
				Node: plugin.Node{
					Unschedulable: true,
				},
				Task: plugin.Task{
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
		if status.Error == "" {
			fmt.Println(item.name + ": pass")
		} else {
			fmt.Println(status.Error)
		}
	}
}

func helper(path, name string, args *plugin.Args) (common.Status, error) {
	config := gop.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "plugin-filter",
		MagicCookieValue: "plugin-filter",
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin-filter",
		Output: os.Stderr,
		Level:  hclog.Error,
	})

	plugins := map[string]gop.Plugin{
		name: &common.FilterPlugin{},
	}

	client := gop.NewClient(&gop.ClientConfig{
		Cmd:             exec.Command(path),
		HandshakeConfig: config,
		Logger:          logger,
		Plugins:         plugins,
	})
	defer client.Kill()

	rpcClient, _ := client.Client()
	raw, _ := rpcClient.Dispense(name)
	n := raw.(common.Filter)
	status := n.Filter(args)

	return status, nil
}
