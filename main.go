package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	gop "github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"

	"github.com/pipego/scheduler/common"
	"github.com/pipego/scheduler/plugin"
)

type config struct {
	args *common.Args
	name string
	path string
}

var (
	configs = []config{
		// Plugin: NodeName
		{
			args: &common.Args{
				Node: common.Node{
					Name: "Node",
				},
				Task: common.Task{
					NodeName: "Node",
				},
			},
			name: "NodeName",
			path: "./plugin/filter-nodename",
		},
		{
			args: &common.Args{
				Node: common.Node{
					Name: "Node",
				},
				Task: common.Task{
					NodeName: "Task",
				},
			},
			name: "NodeName",
			path: "./plugin/filter-nodename",
		},
		// Plugin: NodeResourcesFit
		{
			args: &common.Args{
				Node: common.Node{
					AllocatableResource: common.Resource{
						MilliCPU: 400,
					},
					RequestedResource: common.Resource{
						MilliCPU: 200,
					},
				},
				Task: common.Task{
					RequestedResource: common.Resource{
						MilliCPU: 100,
					},
				},
			},
			name: "NodeResourcesFit",
			path: "./plugin/filter-noderesourcesfit",
		},
		{
			args: &common.Args{
				Node: common.Node{
					AllocatableResource: common.Resource{
						MilliCPU: 400,
					},
					RequestedResource: common.Resource{
						MilliCPU: 200,
					},
				},
				Task: common.Task{
					RequestedResource: common.Resource{
						MilliCPU: 500,
					},
				},
			},
			name: "NodeResourcesFit",
			path: "./plugin/filter-noderesourcesfit",
		},
		// Plugin: NodeAffinity
		{
			args: &common.Args{
				Node: common.Node{
					Label: "ssd",
				},
				Task: common.Task{
					NodeSelectors: []string{"ssd"},
				},
			},
			name: "NodeAffinity",
			path: "./plugin/filter-nodeaffinity",
		},
		{
			args: &common.Args{
				Node: common.Node{
					Label: "ssd",
				},
				Task: common.Task{
					NodeSelectors: []string{"hdd"},
				},
			},
			name: "NodeAffinity",
			path: "./plugin/filter-nodeaffinity",
		},
		// Plugin: NodeUnschedulable
		{
			args: &common.Args{
				Node: common.Node{
					Unschedulable: true,
				},
				Task: common.Task{
					ToleratesUnschedulable: true,
				},
			},
			name: "NodeUnschedulable",
			path: "./plugin/filter-nodeunschedulable",
		},
		{
			args: &common.Args{
				Node: common.Node{
					Unschedulable: true,
				},
				Task: common.Task{
					ToleratesUnschedulable: false,
				},
			},
			name: "NodeUnschedulable",
			path: "./plugin/filter-nodeunschedulable",
		},
	}

	handshake = gop.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "plugin",
		MagicCookieValue: "plugin",
	}

	logger = hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stderr,
		Level:  hclog.Error,
	})
)

func main() {
	for _, item := range configs {
		p, _ := filepath.Abs(item.path)
		if status, err := helper(p, item.name, item.args); err == nil {
			if status.Error == "" {
				fmt.Println(item.name + ": pass")
			} else {
				fmt.Println(status.Error)
			}
		} else {
			fmt.Println(err.Error())
		}
	}
}

func helper(path, name string, args *common.Args) (plugin.FilterResult, error) {
	plugins := map[string]gop.Plugin{
		name: &plugin.Filter{},
	}

	client := gop.NewClient(&gop.ClientConfig{
		Cmd:             exec.Command(path),
		HandshakeConfig: handshake,
		Logger:          logger,
		Plugins:         plugins,
	})
	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		return plugin.FilterResult{}, errors.Wrap(err, "failed to init client")
	}

	raw, err := rpcClient.Dispense(name)
	if err != nil {
		return plugin.FilterResult{}, errors.Wrap(err, "failed to dispense instance")
	}

	n := raw.(plugin.FilterImpl)
	status := n.Run(args)

	return status, nil
}
