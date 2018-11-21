package load_balance

import (
	proto "../common_proto"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"strconv"
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
	m_isConnect        bool
	m_conn             *zk.Conn
	m_connChan         chan bool
}

func (this *CZkAdapter) init(conns *[]proto.CConnectProperty, pathPrefix string, connTimeoutS int) (<-chan bool, error) {
	this.m_pathPrefix = pathPrefix
	this.ZkBaseInit(conns, connTimeoutS, this)
	this.m_connChan = make(chan bool, 1)
	go this.sync()
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
	this.m_isConnect = true
	this.m_connChan <- true
	close(this.m_connChan)
	this.m_connChan = make(chan bool, 1)
	return nil
}

func (this *CZkAdapter) EventCallback(event zk.Event) {
	if event.Type == zk.EventNodeCreated || event.Type == zk.EventNodeDeleted || event.Type == zk.EventNodeDataChanged {
		fmt.Println("[INFO] node change")
	}
	if event.State == zk.StateDisconnected {
		this.m_isConnect = false
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
		if this.m_isConnect == false {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		fmt.Println("[INFO] sync data start")
		this.syncData()
		fmt.Println("[INFO] sync data end")
		time.Sleep(3 * time.Second)
	}
}

func (*CZkAdapter) nodeData2hash(data *proto.CNodeData) int {
	str := strings.Join([]string{data.ServerIp, strconv.FormatInt(int64(data.ServerPort), 10), data.ServerUniqueCode, strconv.FormatInt(int64(data.Weight), 10)}, ":")
	return toHash([]byte(str))
}

func (this *CZkAdapter) normalNodeIsUpdate(news *[]proto.CNodeData, olds *[]proto.CNodeData) (_adds *[]proto.CNodeData, _dels *[]proto.CNodeData) {
	var adds []proto.CNodeData
	var deletes []proto.CNodeData
	var oldMap sync.Map
	var oldMapCopy sync.Map
	for _, old := range *olds {
		hash := this.nodeData2hash(&old)
		oldMap.Store(hash, &old)
		oldMapCopy.Store(hash, &old)
	}
	for _, _new := range *news {
		hash := this.nodeData2hash(&_new)
		_, ok := oldMap.Load(hash)
		if ok {
			// both exist -> no change
		} else {
			// new exist, old is not exist -> add or update
			// if is update -> delete first, then add
			this.deleteNodeData(&_new, &adds)
			adds = append(adds, _new)
		}
		oldMapCopy.Delete(hash)
	}
	f := func(k, v interface{}) bool {
		d := v.(*proto.CNodeData)
		deletes = append(deletes, *d)
		return true
	}
	oldMapCopy.Range(f)
	return &adds, &deletes
}

func (this *CZkAdapter) deleteNodeData(delData *proto.CNodeData, datas *[]proto.CNodeData) {
	delHash := this.nodeData2hash(delData)
	for i, d := range *datas {
		hash := this.nodeData2hash(&d)
		if hash == delHash {
			*datas = append((*datas)[:i], (*datas)[i+1:]...)
			break
		}
	}
}

func (this *CZkAdapter) syncData() error {
	path := this.JoinPathPrefix(&this.m_pathPrefix, nil)
	servers, _, err := this.m_conn.Children(*path)
	if err != nil {
		fmt.Println("[ERROR] get server error, path: ", *path)
		return err
	}
	var serverDataCopy sync.Map
	f1 := func(k, v interface{}) bool {
		serverDataCopy.Store(k, v)
		return true
	}
	this.m_serverData.Range(f1)
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
		serverInfo, ok := this.m_serverData.Load(server)
		if !ok {
			this.m_serverData.Store(server, CDataItem{masterNode: &masterNode, normalNodes: &normalNodes, isChanged: false})
		} else {
			info := serverInfo.(CDataItem)
			adds, deletes := this.normalNodeIsUpdate(&normalNodes, info.normalNodes)
			addLen := len(*adds)
			delLen := len(*deletes)
			if addLen > 0 || delLen > 0 {
				fmt.Println("[INFO] normal node changed")
				for _, add := range *adds {
					*info.normalNodes = append(*info.normalNodes, add)
				}
				for _, del := range *deletes {
					this.deleteNodeData(&del, info.normalNodes)
				}
				info.isChanged = true
			} else {
				info.isChanged = false
			}
			this.m_serverData.Store(server, info)
		}
		// remove new have from old -> remain: need delete
		serverDataCopy.Delete(server)
	}
	// delete from serverMap
	f := func(k, v interface{}) bool {
		fmt.Println("[INFO] delete server: ", k.(string))
		this.m_serverData.Delete(k.(string))
		return true
	}
	serverDataCopy.Range(f)
	return nil
}
