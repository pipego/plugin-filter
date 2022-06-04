package main

import (
	"strings"

	gop "github.com/hashicorp/go-plugin"

	"github.com/pipego/scheduler/common"
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

func (n *NodeResourcesFit) Run(args *common.Args) plugin.FilterResult {
	var status plugin.FilterResult

	if args.Task.RequestedResource.MilliCPU <= 0 &&
		args.Task.RequestedResource.Memory <= 0 &&
		args.Task.RequestedResource.Storage <= 0 {
		status.Error = ErrReasonResourcesFit + ": " + "invalid resource"
		return status
	}

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

func (n *NodeResourcesFit) fit(task *common.Task, node *common.Node) []InsufficientResource {
	insufficientResources := make([]InsufficientResource, 0, 4)

	if task.RequestedResource.MilliCPU > (node.AllocatableResource.MilliCPU - node.RequestedResource.MilliCPU) {
		insufficientResources = append(insufficientResources, InsufficientResource{
			ResourceName: common.ResourceCPU,
			Reason:       "insufficient cpu",
			Requested:    task.RequestedResource.MilliCPU,
			Used:         node.RequestedResource.MilliCPU,
			Capacity:     node.AllocatableResource.MilliCPU,
		})
	}

	if task.RequestedResource.Memory > (node.AllocatableResource.Memory - node.RequestedResource.Memory) {
		insufficientResources = append(insufficientResources, InsufficientResource{
			ResourceName: common.ResourceMemory,
			Reason:       "insufficient memory",
			Requested:    task.RequestedResource.Memory,
			Used:         node.RequestedResource.Memory,
			Capacity:     node.AllocatableResource.Memory,
		})
	}

	if task.RequestedResource.Storage > (node.AllocatableResource.Storage - node.RequestedResource.Storage) {
		insufficientResources = append(insufficientResources, InsufficientResource{
			ResourceName: common.ResourceStorage,
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
		MagicCookieKey:   "plugin",
		MagicCookieValue: "plugin",
	}

	pluginMap := map[string]gop.Plugin{
		"NodeResourcesFit": &plugin.Filter{Impl: &NodeResourcesFit{}},
	}

	gop.Serve(&gop.ServeConfig{
		HandshakeConfig: config,
		Plugins:         pluginMap,
	})
}
