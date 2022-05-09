package main

import (
	gop "github.com/hashicorp/go-plugin"

	"github.com/pipego/plugin-filter/common"
	"github.com/pipego/scheduler/plugin"
)

const (
	ErrReasonName = "NodeName: node(s) didn't match the requested node name"
)

type NodeName struct{}

func (n *NodeName) Filter(args *plugin.Args) common.Status {
	var status common.Status

	if args.Task.NodeName != "" && args.Task.NodeName != args.Node.Name {
		status.Error = ErrReasonName
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

	pluginMap := map[string]gop.Plugin{
		"NodeName": &common.FilterPlugin{Impl: &NodeName{}},
	}

	gop.Serve(&gop.ServeConfig{
		HandshakeConfig: config,
		Plugins:         pluginMap,
	})
}
