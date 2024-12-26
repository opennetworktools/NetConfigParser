package parser

import (
	"github.com/opennetworktools/NetConfigParser/internal/net"
)

type Configs struct {
	Metadata      net.Metadata
	BGP           net.BGP
	Interfaces    net.Interfaces
	RouteMaps     []net.RouteMap
	IPAccessLists []net.AccessList
	IPPrefixLists []net.PrefixList
}

type ConfigParser interface {
	GetConfigs() *Configs
	PrintParserType()
	ParseMetadata(string) net.Metadata
	ParseConfig() error
	ParseBGPBlock(string) net.BGP
	ParseInterfacesBlock([]string) net.Interfaces
	ParseInterfaceString(string, string) net.Interface
	ParseRouteMapBlock([]string) []net.RouteMap
	ParseIPAccessListBlock([]string) []net.AccessList
	ParseIPPrefixListBlock([]string) []net.PrefixList
}

func NewIOSXEParser(path string) *IOSXEParser {
	return &IOSXEParser{
		FilePath: path,
		Features: make([]string, 0),
		Configs:  Configs{},
	}
}
