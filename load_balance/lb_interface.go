package load_balance

import (
	proto "../common_proto"
	"sync"
)

var (
	ServerModeZookeeper     = "zookeeper"
	ServerModeZookeeperHttp = "zookeeper_http"
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
	Get(serverName string, extraData interface{}) (*proto.CNodeData, error)
	init(loadBlance ILoadBlance) error
}

type ICallback interface {
	MasterNodeChange(data *proto.CNodeData, userData interface{})
	NormalNodeChange(data *proto.CNodeData, userData interface{})
	ServerBeDeleted(serverName *string, userData interface{})
}

type CDataItem struct {
	masterNode  *proto.CNodeData
	normalNodes *[]proto.CNodeData
	isChanged   bool
}

type ILoadBlance interface {
	SetCallback(callback ICallback, userData interface{})
	GetMasterNode(serverName string) (*proto.CNodeData, error)
	GetNormalNodeAlgorithm(algorithm string) INormalNodeAlgorithm
	init(conns *[]proto.CConnectProperty, pathPrefix string, connTimeoutS int) (<-chan bool, error)
	findAllServerData() (*sync.Map, error)
	findServerData(serverName string) (*CDataItem, error)
	nodeData2hash(data *proto.CNodeData) int
}

func New(serverMode string, conns *[]proto.CConnectProperty, pathPrefix string, connTimeoutS int) (ILoadBlance, <-chan bool) {
	if serverMode == ServerModeZookeeper {
		adapter := &CZkAdapter{}
		ch, err := adapter.init(conns, pathPrefix, connTimeoutS)
		if err != nil {
			return nil, nil
		}
		return adapter, ch
	} else if serverMode == ServerModeZookeeperHttp {
		adapter := &CZkHttpAdapter{}
		ch, err := adapter.init(conns, pathPrefix, connTimeoutS)
		if err != nil {
			return nil, nil
		}
		return adapter, ch
	}
	return nil, nil
}
