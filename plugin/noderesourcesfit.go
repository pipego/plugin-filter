package main

import (
	"strings"

	gop "github.com/hashicorp/go-plugin"

	"github.com/pipego/plugin-filter/common"
	"github.com/pipego/scheduler/plugin"
)

const (
	ErrReasonResourcesFit = "NodeResourcesFit: node(s) didn't fit the resources"
)

type NodeResourcesFit struct{}

type InsufficientResource struct {
	ResourceName string
	Reason       string
	Requested    int64
	Used         int64
	Capacity     int64
}

func (n *NodeResourcesFit) Run(args *plugin.Args) plugin.FilterResult {
	var status plugin.FilterResult

	insufficientResources := n.fit(&args.Task, &args.Node)
	if len(insufficientResources) != 0 {
		failureReasons := make([]string, 0, len(insufficientResources))
		for _, item := range insufficientResources {
			failureReasons = append(failureReasons, item.Reason)
		}
		status.Error = ErrReasonResourcesFit + ": " + strings.Join(failureReasons, ", ")
	}

	return status
}

func (n *NodeResourcesFit) fit(task *plugin.Task, node *plugin.Node) []InsufficientResource {
	insufficientResources := make([]InsufficientResource, 0, 4)

	if task.RequestedResource.MilliCPU == 0 &&
		task.RequestedResource.Memory == 0 &&
		task.RequestedResource.Storage == 0 {
		return insufficientResources
	}

	if task.RequestedResource.MilliCPU > (node.AllocatableResource.MilliCPU - node.RequestedResource.MilliCPU) {
		insufficientResources = append(insufficientResources, InsufficientResource{
			ResourceName: plugin.ResourceCPU,
			Reason:       "insufficient cpu",
			Requested:    task.RequestedResource.MilliCPU,
			Used:         node.RequestedResource.MilliCPU,
			Capacity:     node.AllocatableResource.MilliCPU,
		})
	}

	if task.RequestedResource.Memory > (node.AllocatableResource.Memory - node.RequestedResource.Memory) {
		insufficientResources = append(insufficientResources, InsufficientResource{
			ResourceName: plugin.ResourceMemory,
			Reason:       "insufficient memory",
			Requested:    task.RequestedResource.Memory,
			Used:         node.RequestedResource.Memory,
			Capacity:     node.AllocatableResource.Memory,
		})
	}

	if task.RequestedResource.Storage > (node.AllocatableResource.Storage - node.RequestedResource.Storage) {
		insufficientResources = append(insufficientResources, InsufficientResource{
			ResourceName: plugin.ResourceStorage,
			Reason:       "insufficient storage",
			Requested:    task.RequestedResource.Storage,
			Used:         node.RequestedResource.Storage,
			Capacity:     node.AllocatableResource.Storage,
		})
	}

	return insufficientResources
}

// nolint:typecheck
func main() {
	config := gop.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "plugin-filter",
		MagicCookieValue: "plugin-filter",
	}

	pluginMap := map[string]gop.Plugin{
		"NodeResourcesFit": &common.FilterPlugin{Impl: &NodeResourcesFit{}},
	}

	gop.Serve(&gop.ServeConfig{
		HandshakeConfig: config,
		Plugins:         pluginMap,
	})
}
