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

type CZkDataItem struct {
	nodeType string
	nodeData proto.CNodeData
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
	fmt.Println(this.m_serverData)
	path := this.JoinPathPrefix(&this.m_pathPrefix, &serverName)
	childrens, _, err := this.m_conn.Children(*path)
	if err != nil {
		fmt.Println("[ERROR] get children error")
		return nil, err
	}
	isFind := false
	var masterNode string = ""
	for _, child := range childrens {
		if child == proto.MasterNode {
			isFind = true
			childPath := strings.Join([]string{*path, child}, "/")
			b, _, err := this.m_conn.Get(childPath)
			if err != nil {
				fmt.Println("[ERROR] get nodedata error, path: ", childPath)
				return nil, err
			}
			masterNode = string(b)
			break
		}
	}
	if isFind == false {
		return nil, errors.New("is not exist")
	}
	fmt.Println("[DEBUG] nodedata: ", masterNode)
	data := proto.CNodeData{}
	err = json.Unmarshal([]byte(masterNode), &data)
	if err != nil {
		fmt.Println("[ERROR] decoder masternode json error")
		return nil, err
	}
	return &data, nil
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
		var items []CZkDataItem
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
			var t string = proto.NormalNode
			if node == proto.MasterNode {
				t = proto.NormalNode
			}
			item := CZkDataItem{nodeType: t, nodeData: data}
			items = append(items, item)
		}
		this.m_serverData.Store(server, items)
	}
	return nil
}
