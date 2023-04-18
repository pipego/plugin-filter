//nolint:typecheck
package main

import (
	gop "github.com/hashicorp/go-plugin"

	"github.com/pipego/scheduler/common"
	"github.com/pipego/scheduler/plugin"
)

const (
	ErrReasonUnschedulable = "NodeUnschedulable: node(s) were unschedulable"
)

type NodeUnschedulable struct{}

func (n *NodeUnschedulable) Run(args *common.Args) plugin.FilterResult {
	var status plugin.FilterResult

	if args.Node.Unschedulable && !args.Task.ToleratesUnschedulable {
		status.Error = ErrReasonUnschedulable
	}

	return status
}

func main() {
	config := gop.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "plugin",
		MagicCookieValue: "plugin",
	}

	var pluginMap = map[string]gop.Plugin{
		"NodeUnschedulable": &plugin.Filter{Impl: &NodeUnschedulable{}},
	}

	gop.Serve(&gop.ServeConfig{
		HandshakeConfig: config,
		Plugins:         pluginMap,
	})
}
