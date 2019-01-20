package load_balance

import (
	proto "../common_proto"
	"sync"
)

var (
	ServerModeZookeeper     = "zookeeper"
	ServerModeZookeeperHttp = "zookeeper_http"
	ServerModeZookeeperMqtt = "zookeeper_mqtt"
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

type CNetInfo struct {
	Host     string
	Port     int
	UserName string
	UserPwd  string
}

type ILoadBlance interface {
	SetCallback(callback ICallback, userData interface{})
	SetTransmitTimeoutS(s int)
	AddRecvNetInfo(topic *string, info *CNetInfo)
	SetNormalNodeAlgorithm(algorithm string) error
	GetMasterNode(serverName string) (*proto.CNodeData, error)
	GetNormalNodeAlgorithm(algorithm string) INormalNodeAlgorithm
	Run(data interface{}) error
	init(conns *[]proto.CConnectProperty, pathPrefix string, connTimeoutS int) (<-chan bool, error)
	findAllServerData() (*sync.Map, error)
	findServerData(serverName string) (*proto.CDataItem, error)
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
	} else if serverMode == ServerModeZookeeperMqtt {
		adapter := &CZkMqttAdapter{}
		ch, err := adapter.init(conns, pathPrefix, connTimeoutS)
		if err != nil {
			return nil, nil
		}
		return adapter, ch
	}
	return nil, nil
}
