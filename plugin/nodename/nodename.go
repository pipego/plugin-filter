//nolint:typecheck
package main

import (
	gop "github.com/hashicorp/go-plugin"

	"github.com/pipego/scheduler/common"
	"github.com/pipego/scheduler/plugin"
)

const (
	ErrReasonName = "NodeName: node(s) didn't match the requested node name"
)

type NodeName struct{}

func (n *NodeName) Run(args *common.Args) plugin.FilterResult {
	var status plugin.FilterResult

	if args.Task.NodeName == "" || args.Task.NodeName != args.Node.Name {
		status.Error = ErrReasonName
	}

	return status
}

func main() {
	config := gop.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "plugin",
		MagicCookieValue: "plugin",
	}

	pluginMap := map[string]gop.Plugin{
		"NodeName": &plugin.Filter{Impl: &NodeName{}},
	}

	gop.Serve(&gop.ServeConfig{
		HandshakeConfig: config,
		Plugins:         pluginMap,
	})
}
