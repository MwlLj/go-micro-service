package load_balance

import (
	proto "../common_proto"
)

var (
	ServerModeZookeeper = "zookeeper"
)

type ICallback interface {
	MasterNodeChange(data *proto.CNodeData, userData interface{})
	NormalNodeChange(data *proto.CNodeData, userData interface{})
}

type ILoadBlance interface {
	SetCallback(callback ICallback, userData interface{})
	GetMasterNode(serverName string) *string
	RoundRobin(serverName string) *string
	WeightRoundRobin(serverName string) *string
	Random(serverName string) *string
	WeightRandom(serverName string) *string
	IpHash(serverName string) *string
	UrlHash(serverName string) *string
	LeastConnections(serverName string) *string
	init(conns *[]proto.CConnectProperty, connTimeoutS int) error
}

func New(serverMode string, conns *[]proto.CConnectProperty, connTimeoutS int) ILoadBlance {
	if serverMode == ServerModeZookeeper {
		adapter := &CZkAdapter{}
		adapter.init(conns, connTimeoutS)
		return adapter
	}
	return nil
}
