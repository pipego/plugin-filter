package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/pipego/plugin-filter/proto"
	"sort"
)

const (
	ErrReasonAffinity = "NodeAffinity: node(s) didn't match Task's node affinity/selector"
)

type NodeAffinity struct{}

func (n *NodeAffinity) Filter(args *proto.Args) proto.Status {
	var status proto.Status
	found := false

	for key, val := range args.Task.NodeSelector {
		if _, ok := args.Node.Label[key]; ok {
			if n.match(args.Node.Label[key], val) {
				found = true
			}
		}
	}

	if !found {
		status.Error = ErrReasonAffinity
	}

	return status
}

func (n *NodeAffinity) match(name string, list []string) bool {
	sort.Strings(list)

	index := sort.SearchStrings(list, name)
	if index >= len(list) || list[index] != name {
		return false
	}

	return true
}

// nolint:typecheck
func main() {
	config := plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "plugin-filter",
		MagicCookieValue: "plugin-filter",
	}

	pluginMap := map[string]plugin.Plugin{
		"NodeAffinity": &proto.FilterPlugin{Impl: &NodeAffinity{}},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: config,
		Plugins:         pluginMap,
	})
}
