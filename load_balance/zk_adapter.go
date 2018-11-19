package load_balance

import (
	proto "../common_proto"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
	"sync"
	"time"
)

var _ = fmt.Println
var _ = errors.New("")
var _ = strings.Join
var _ = json.Unmarshal

type CZkAdapter struct {
	proto.CZkBase
	m_pathPrefix       string
	m_callback         ICallback
	m_callbackUserData interface{}
	m_serverData       sync.Map
	m_conn             *zk.Conn
	m_connChan         chan bool
}

func (this *CZkAdapter) init(conns *[]proto.CConnectProperty, pathPrefix string, connTimeoutS int) (<-chan bool, error) {
	this.m_pathPrefix = pathPrefix
	this.ZkBaseInit(conns, connTimeoutS, this)
	this.m_connChan = make(chan bool, 1)
	return this.m_connChan, nil
}

func (this *CZkAdapter) SetCallback(callback ICallback, userData interface{}) {
	this.m_callback = callback
	this.m_callbackUserData = userData
}

func (this *CZkAdapter) AddConnProperty(conn *proto.CConnectProperty) error {
	return this.AddConnProperty(conn)
}

func (this *CZkAdapter) UpdateConnProperty(conn *proto.CConnectProperty) error {
	return this.UpdateConnProperty(conn)
}

func (this *CZkAdapter) DeleteConnProperty(serviceId *string) error {
	return this.DeleteConnProperty(serviceId)
}

func (this *CZkAdapter) AfterConnect(conn *zk.Conn) error {
	this.m_conn = conn
	this.syncData()
	this.m_connChan <- true
	close(this.m_connChan)
	return nil
}

func (this *CZkAdapter) EventCallback(event zk.Event) {
	if event.Type == zk.EventNodeCreated || event.Type == zk.EventNodeDeleted || event.Type == zk.EventNodeDataChanged {
		fmt.Println("[INFO] node change")
	}
}

func (this *CZkAdapter) GetNormalNodeAlgorithm(algorithm string) INormalNodeAlgorithm {
	if algorithm == AlgorithmRoundRobin {
		return &CRoundRobin{m_loadBlance: this}
	} else if algorithm == AlgorithmWeightRoundRobin {
		return &CWeightRoundRobin{m_loadBlance: this}
	} else if algorithm == AlgorithmRandom {
		return &CRandom{m_loadBlance: this}
	} else if algorithm == AlgorithmWeightRandom {
		return &CWeightRandom{m_loadBlance: this}
	} else if algorithm == AlgorithmIpHash {
		return &CIpHash{m_loadBlance: this}
	} else if algorithm == AlgorithmUrlHash {
		return &CUrlHash{m_loadBlance: this}
	} else if algorithm == AlgorithmLeastConnect {
		return &CLeastConnections{m_loadBlance: this}
	} else {
		return nil
	}
}

func (this *CZkAdapter) GetMasterNode(serverName string) (*proto.CNodeData, error) {
	item, err := this.findServerData(serverName)
	if err != nil {
		fmt.Println("[ERROR] find serverdata error")
		return nil, err
	}
	if item.masterNode == nil {
		return nil, errors.New("is not exist")
	}
	return item.masterNode, nil
}

func (this *CZkAdapter) findServerData(serverName string) (*CDataItem, error) {
	var resultValue interface{} = nil
	f := func(k, v interface{}) bool {
		key := k.(string)
		if key == serverName {
			resultValue = v
			return false
		}
		return true
	}
	this.m_serverData.Range(f)
	if resultValue == nil {
		return nil, errors.New("server not found")
	}
	r := resultValue.(CDataItem)
	return &r, nil
}

func (this *CZkAdapter) sync() error {
	for {
		this.syncData()
		time.Sleep(30 * time.Second)
	}
}

func (this *CZkAdapter) syncData() error {
	path := this.JoinPathPrefix(&this.m_pathPrefix, nil)
	servers, _, err := this.m_conn.Children(*path)
	if err != nil {
		fmt.Println("[ERROR] get server error, path: ", *path)
		return err
	}
	for _, server := range servers {
		// fmt.Println("[DEBUG] server: ", server)
		nodeRootPath := this.JoinPathPrefix(&this.m_pathPrefix, &server)
		nodes, _, err := this.m_conn.Children(*nodeRootPath)
		if err != nil {
			fmt.Println("[WARNING] get node error, path: ", nodeRootPath)
			continue
		}
		var masterNode proto.CNodeData
		var normalNodes []proto.CNodeData
		for _, node := range nodes {
			// fmt.Println("[DEBUG] node: ", node)
			nodePath := strings.Join([]string{*path, server, node}, "/")
			// fmt.Println("[DEBUG] nodePath: ", nodePath)
			b, _, err := this.m_conn.Get(nodePath)
			if err != nil {
				fmt.Println("[WARNING] get nodedata error, path: ", nodePath)
				continue
			}
			// fmt.Println("[DEBUG] nodeData: ", string(b))
			data := proto.CNodeData{}
			err = json.Unmarshal(b, &data)
			if err != nil {
				fmt.Println("[WARNING] decoder nodedata error, node")
				continue
			}
			if node == proto.MasterNode {
				masterNode = data
			} else {
				normalNodes = append(normalNodes, data)
			}
		}
		this.m_serverData.Store(server, CDataItem{masterNode: &masterNode, normalNodes: &normalNodes})
	}
	return nil
}
