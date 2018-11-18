package load_balance

import (
	proto "../common_proto"
)

var (
	ServerModeNocacheZookeeper   = "nocache_zookeeper"
	ServerModeWithcacheZookeeper = "withcache_zookeeper"
)

type ICallback interface {
	MasterNodeChange(data *proto.CNodeData, userData interface{})
	NormalNodeChange(data *proto.CNodeData, userData interface{})
}

type ILoadBlance interface {
	SetCallback(callback ICallback, userData interface{})
	GetMasterNode(serverName string) (*proto.CNodeData, error)
	RoundRobin(serverName string) (*proto.CNodeData, error)
	WeightRoundRobin(serverName string) (*proto.CNodeData, error)
	Random(serverName string) (*proto.CNodeData, error)
	WeightRandom(serverName string) (*proto.CNodeData, error)
	IpHash(serverName string) (*proto.CNodeData, error)
	UrlHash(serverName string) (*proto.CNodeData, error)
	LeastConnections(serverName string) (*proto.CNodeData, error)
	init(conns *[]proto.CConnectProperty, pathPrefix string, connTimeoutS int) error
}

func New(serverMode string, conns *[]proto.CConnectProperty, pathPrefix string, connTimeoutS int) ILoadBlance {
	if serverMode == ServerModeNocacheZookeeper {
		adapter := &CNocacheZkAdapter{}
		adapter.init(conns, pathPrefix, connTimeoutS)
		return adapter
	} else if serverMode == ServerModeWithcacheZookeeper {
		adapter := &CWithcacheZkAdapter{}
		adapter.init(conns, pathPrefix, connTimeoutS)
		return adapter
	}
	return nil
}
