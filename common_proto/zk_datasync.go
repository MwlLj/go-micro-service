package common_proto

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"hash/crc32"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CDataItem struct {
	MasterNode  *CNodeData
	NormalNodes *[]CNodeData
	IsChanged   bool
}

type CZkDataSync struct {
	CZkBase
	m_pathPrefix string
	m_serverData sync.Map
	m_isConnect  bool
	m_conn       *zk.Conn
	m_connChan   chan bool
}

func toHash(b []byte) int {
	v := int(crc32.ChecksumIEEE(b))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

func (this *CZkDataSync) Init(conns *[]CConnectProperty, pathPrefix string, connTimeoutS int) (<-chan bool, error) {
	this.m_pathPrefix = pathPrefix
	this.ZkBaseInit(conns, connTimeoutS, this)
	this.m_connChan = make(chan bool, 1)
	go this.sync()
	return this.m_connChan, nil
}

func (this *CZkDataSync) AfterConnect(conn *zk.Conn) error {
	this.m_conn = conn
	this.syncData()
	this.m_isConnect = true
	this.m_connChan <- true
	close(this.m_connChan)
	this.m_connChan = make(chan bool, 1)
	return nil
}

func (this *CZkDataSync) EventCallback(event zk.Event) {
	if event.Type == zk.EventNodeCreated || event.Type == zk.EventNodeDeleted || event.Type == zk.EventNodeDataChanged {
		fmt.Println("[INFO] node change")
	}
	if event.State == zk.StateDisconnected {
		this.m_isConnect = false
	}
}

func (this *CZkDataSync) FindAllServerData() (*sync.Map, error) {
	return &this.m_serverData, nil
}

func (this *CZkDataSync) FindServerData(serverName string) (*CDataItem, error) {
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

func (this *CZkDataSync) sync() error {
	for {
		if this.m_isConnect == false {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		fmt.Println("[INFO] sync data start")
		this.syncData()
		fmt.Println("[INFO] sync data end")
		time.Sleep(30 * time.Second)
	}
}

func (*CZkDataSync) NodeData2hash(data *CNodeData) int {
	str := strings.Join([]string{data.ServerIp, strconv.FormatInt(int64(data.ServerPort), 10), data.ServerUniqueCode, strconv.FormatInt(int64(data.Weight), 10)}, ":")
	return toHash([]byte(str))
}

func (this *CZkDataSync) normalNodeIsUpdate(news *[]CNodeData, olds *[]CNodeData) (_adds *[]CNodeData, _dels *[]CNodeData) {
	var adds []CNodeData
	var deletes []CNodeData
	var oldMap sync.Map
	var oldMapCopy sync.Map
	for _, old := range *olds {
		oldTmp := old
		hash := this.NodeData2hash(&oldTmp)
		oldMap.Store(hash, &oldTmp)
		oldMapCopy.Store(hash, &oldTmp)
	}
	for _, _new := range *news {
		newTmp := _new
		hash := this.NodeData2hash(&newTmp)
		_, ok := oldMap.Load(hash)
		if ok {
			// both exist -> no change
		} else {
			// new exist, old is not exist -> add or update
			// if is update -> delete first, then add
			this.deleteNodeData(&newTmp, &adds)
			adds = append(adds, newTmp)
		}
		oldMapCopy.Delete(hash)
	}
	f := func(k, v interface{}) bool {
		d := v.(*CNodeData)
		deletes = append(deletes, *d)
		return true
	}
	oldMapCopy.Range(f)
	return &adds, &deletes
}

func (this *CZkDataSync) deleteNodeData(delData *CNodeData, datas *[]CNodeData) {
	delHash := this.NodeData2hash(delData)
	for i, d := range *datas {
		tmp := d
		hash := this.NodeData2hash(&tmp)
		if hash == delHash {
			*datas = append((*datas)[:i], (*datas)[i+1:]...)
			break
		}
	}
}

func (this *CZkDataSync) syncData() error {
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
		serverTmp := server
		nodeRootPath := this.JoinPathPrefix(&this.m_pathPrefix, &serverTmp)
		nodes, _, err := this.m_conn.Children(*nodeRootPath)
		if err != nil {
			fmt.Println("[WARNING] get node error, path: ", nodeRootPath)
			continue
		}
		var masterNode CNodeData
		var normalNodes []CNodeData
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
			data := CNodeData{}
			err = json.Unmarshal(b, &data)
			if err != nil {
				fmt.Println("[WARNING] decoder nodedata error, node")
				continue
			}
			if node == MasterNode {
				masterNode = data
			} else {
				normalNodes = append(normalNodes, data)
			}
		}
		serverInfo, ok := this.m_serverData.Load(server)
		if !ok {
			this.m_serverData.Store(server, CDataItem{MasterNode: &masterNode, NormalNodes: &normalNodes, IsChanged: false})
		} else {
			info := serverInfo.(CDataItem)
			adds, deletes := this.normalNodeIsUpdate(&normalNodes, info.NormalNodes)
			addLen := len(*adds)
			delLen := len(*deletes)
			if addLen > 0 || delLen > 0 {
				fmt.Println("[INFO] normal node changed")
				for _, add := range *adds {
					*info.NormalNodes = append(*info.NormalNodes, add)
				}
				for _, del := range *deletes {
					tmp := del
					this.deleteNodeData(&tmp, info.NormalNodes)
				}
				info.IsChanged = true
			} else {
				info.IsChanged = false
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
