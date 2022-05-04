package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/pipego/plugin-filter/proto"
)

const (
	ErrReasonSelector = "NodeSelector: node(s) didn't match the selector"
)

type NodeSelector struct{}

func (n *NodeSelector) Filter(args *proto.Args) proto.Status {
	var status proto.Status

	// TODO
	status.Error = ErrReasonSelector

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
		"NodeSelector": &proto.FilterPlugin{Impl: &NodeSelector{}},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: config,
		Plugins:         pluginMap,
	})
}
