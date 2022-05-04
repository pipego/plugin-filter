package main

import (
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/pipego/plugin-filter/proto"
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

func (n *NodeResourcesFit) Filter(args *proto.Args) proto.Status {
	var status proto.Status

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

func (n *NodeResourcesFit) fit(task *proto.Task, node *proto.Node) []InsufficientResource {
	insufficientResources := make([]InsufficientResource, 0, 4)

	if task.RequestedResource.MilliCPU == 0 &&
		task.RequestedResource.Memory == 0 &&
		task.RequestedResource.Storage == 0 {
		return insufficientResources
	}

	if task.RequestedResource.MilliCPU > (node.AllocatableResource.MilliCPU - node.RequestedResource.MilliCPU) {
		insufficientResources = append(insufficientResources, InsufficientResource{
			ResourceName: proto.ResourceCPU,
			Reason:       "insufficient cpu",
			Requested:    task.RequestedResource.MilliCPU,
			Used:         node.RequestedResource.MilliCPU,
			Capacity:     node.AllocatableResource.MilliCPU,
		})
	}

	if task.RequestedResource.Memory > (node.AllocatableResource.Memory - node.RequestedResource.Memory) {
		insufficientResources = append(insufficientResources, InsufficientResource{
			ResourceName: proto.ResourceMemory,
			Reason:       "insufficient memory",
			Requested:    task.RequestedResource.Memory,
			Used:         node.RequestedResource.Memory,
			Capacity:     node.AllocatableResource.Memory,
		})
	}

	if task.RequestedResource.Storage > (node.AllocatableResource.Storage - node.RequestedResource.Storage) {
		insufficientResources = append(insufficientResources, InsufficientResource{
			ResourceName: proto.ResourceStorage,
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
	config := plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "plugin-filter",
		MagicCookieValue: "plugin-filter",
	}

	pluginMap := map[string]plugin.Plugin{
		"NodeResourcesFit": &proto.FilterPlugin{Impl: &NodeResourcesFit{}},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: config,
		Plugins:         pluginMap,
	})
}
