package main

import (
	gop "github.com/hashicorp/go-plugin"

	"github.com/pipego/plugin-filter/common"
	"github.com/pipego/scheduler/plugin"
)

const (
	ErrReasonUnschedulable = "NodeUnschedulable: node(s) were unschedulable"
)

type NodeUnschedulable struct{}

func (n *NodeUnschedulable) Filter(args *plugin.Args) common.Status {
	var status common.Status

	if args.Node.Unschedulable && !args.Task.ToleratesUnschedulable {
		status.Error = ErrReasonUnschedulable
	}

	return status
}

// nolint:typecheck
func main() {
	config := gop.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "plugin-filter",
		MagicCookieValue: "plugin-filter",
	}

	var pluginMap = map[string]gop.Plugin{
		"NodeUnschedulable": &common.FilterPlugin{Impl: &NodeUnschedulable{}},
	}

	gop.Serve(&gop.ServeConfig{
		HandshakeConfig: config,
		Plugins:         pluginMap,
	})
}
