//nolint:typecheck
package main

import (
	"sort"

	gop "github.com/hashicorp/go-plugin"

	"github.com/pipego/scheduler/common"
	"github.com/pipego/scheduler/plugin"
)

const (
	ErrReasonAffinity = "NodeAffinity: node(s) didn't match Task's node affinity/selector"
)

type NodeAffinity struct{}

func (n *NodeAffinity) Run(args *common.Args) plugin.FilterResult {
	var status plugin.FilterResult
	found := false

	for _, item := range args.Task.NodeSelectors {
		if args.Node.Label == item {
			found = true
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

func main() {
	config := gop.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "plugin",
		MagicCookieValue: "plugin",
	}

	pluginMap := map[string]gop.Plugin{
		"NodeAffinity": &plugin.Filter{Impl: &NodeAffinity{}},
	}

	gop.Serve(&gop.ServeConfig{
		HandshakeConfig: config,
		Plugins:         pluginMap,
	})
}
