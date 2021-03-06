package main

import (
	plugin "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/iptb-plugins/docker"
	testbedi "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/iptb/testbed/interfaces"
)

var PluginName string
var NewNode testbedi.NewNodeFunc
var GetAttrList testbedi.GetAttrListFunc
var GetAttrDesc testbedi.GetAttrDescFunc

func init() {
	PluginName = plugin.PluginName
	NewNode = plugin.NewNode
	GetAttrList = plugin.GetAttrList
	GetAttrDesc = plugin.GetAttrDesc
}
