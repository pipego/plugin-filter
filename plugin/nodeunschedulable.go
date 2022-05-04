package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/pipego/plugin-filter/proto"
)

const (
	ErrReasonUnschedulable = "node(s) were unschedulable"
)

type NodeUnschedulable struct{}

func (n *NodeUnschedulable) Filter(args *proto.Args) proto.Status {
	var status proto.Status

	if args.Node.Unschedulable && !args.Task.ToleratesUnschedulable {
		status.Error = ErrReasonUnschedulable
	}

	return status
}

// nolint:typecheck
func main() {
	config := plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "plugin-filter",
		MagicCookieValue: "plugin-filter",
	}

	var pluginMap = map[string]plugin.Plugin{
		"NodeUnschedulable": &proto.FilterPlugin{Impl: &NodeUnschedulable{}},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: config,
		Plugins:         pluginMap,
	})
}
