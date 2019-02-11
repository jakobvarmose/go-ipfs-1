package testbed

import (
	"fmt"

	"github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/iptb/testbed/interfaces"
)

// NodeSpec represents a node's specification
type NodeSpec struct {
	Type  string
	Dir   string
	Attrs map[string]string
}

// IptbPlugin contains exported symbols from loaded plugins
type IptbPlugin struct {
	From        string
	NewNode     testbedi.NewNodeFunc
	GetAttrList testbedi.GetAttrListFunc
	GetAttrDesc testbedi.GetAttrDescFunc
	PluginName  string
	BuiltIn     bool
}

var plugins map[string]IptbPlugin

func init() {
	plugins = make(map[string]IptbPlugin)
}

// GetPlugin returns a plugin registered with RegisterPlugin
func GetPlugin(name string) (IptbPlugin, bool) {
	plg, ok := plugins[name]
	return plg, ok
}

// RegisterPlugin registers a plugin, the `force` flag can be passed to
// override any plugin registered under the same IptbPlugin.PluginName
func RegisterPlugin(plg IptbPlugin, force bool) (bool, error) {
	overloaded := false

	if pl, exists := plugins[plg.PluginName]; exists && !force {
		if pl.BuiltIn {
			overloaded = true
		} else {
			return false, fmt.Errorf("plugin %s already loaded from %s", pl.PluginName, pl.From)
		}
	}

	plugins[plg.PluginName] = plg

	return overloaded, nil

}

// LoadPlugin loads a plugin from `path`
func LoadPlugin(path string) (*IptbPlugin, error) {
	return loadPlugin(path)
}

// LoadPluginCore loads core symbols from a golang plugin into an IptbPlugin

// LoadPluginCore loads attr symbols from a golang plugin into an IptbPlugin


func loadPlugin(path string) (*IptbPlugin, error) {
	return nil, nil
}

// Load uses plugins registered with RegisterPlugin to construct a Core node
// from the NodeSpec
func (ns *NodeSpec) Load() (testbedi.Core, error) {
	pluginName := ns.Type

	if plg, ok := plugins[pluginName]; ok {
		return plg.NewNode(ns.Dir, ns.Attrs)
	}

	return nil, fmt.Errorf("Could not find plugin %s", pluginName)
}

// SetAttr sets an attribute on the NodeSpec
func (ns *NodeSpec) SetAttr(attr string, val string) {
	ns.Attrs[attr] = val
}

// GetAttr gets an attribute from the NodeSpec
func (ns *NodeSpec) GetAttr(attr string) (string, error) {
	if v, ok := ns.Attrs[attr]; ok {
		return v, nil
	}

	return "", fmt.Errorf("Attr not set")
}
