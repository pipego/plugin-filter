package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/pipego/plugin-filter/proto"
)

const (
	ErrReasonResourcesFit = "NodeResourcesFit: node(s) didn't fit the resources"
)

type NodeResourcesFit struct{}

func (n *NodeResourcesFit) Filter(args *proto.Args) proto.Status {
	var status proto.Status

	// TODO
	status.Error = ErrReasonResourcesFit

	return status
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
