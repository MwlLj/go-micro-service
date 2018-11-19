package load_balance

import (
	proto "../common_proto"
)

var (
	ServerModeZookeeper = "zookeeper"
)

var (
	AlgorithmRoundRobin       = "roundrobin"
	AlgorithmWeightRoundRobin = "weightroundrobin"
	AlgorithmRandom           = "random"
	AlgorithmWeightRandom     = "weightrandom"
	AlgorithmIpHash           = "iphash"
	AlgorithmUrlHash          = "urlhash"
	AlgorithmLeastConnect     = "leastconnect"
)

type INormalNodeAlgorithm interface {
	Get(serverName string) (*proto.CNodeData, error)
}

type ICallback interface {
	MasterNodeChange(data *proto.CNodeData, userData interface{})
	NormalNodeChange(data *proto.CNodeData, userData interface{})
}

type ILoadBlance interface {
	SetCallback(callback ICallback, userData interface{})
	GetMasterNode(serverName string) (*proto.CNodeData, error)
	GetNormalNodeAlgorithm(algorithm string) INormalNodeAlgorithm
	init(conns *[]proto.CConnectProperty, pathPrefix string, connTimeoutS int) (<-chan bool, error)
}

func New(serverMode string, conns *[]proto.CConnectProperty, pathPrefix string, connTimeoutS int) (ILoadBlance, <-chan bool) {
	if serverMode == ServerModeZookeeper {
		adapter := &CZkAdapter{}
		ch, err := adapter.init(conns, pathPrefix, connTimeoutS)
		if err != nil {
			return nil, nil
		}
		return adapter, ch
	}
	return nil, nil
}
