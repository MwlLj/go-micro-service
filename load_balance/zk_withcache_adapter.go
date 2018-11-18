package load_balance

import (
	proto "../common_proto"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
)

var _ = fmt.Println
var _ = errors.New("")
var _ = strings.Join
var _ = json.Unmarshal

type CWithcacheZkAdapter struct {
	proto.CZkBase
	m_pathPrefix       string
	m_callback         ICallback
	m_callbackUserData interface{}
	m_conn             *zk.Conn
}

func (this *CWithcacheZkAdapter) init(conns *[]proto.CConnectProperty, pathPrefix string, connTimeoutS int) error {
	this.m_pathPrefix = pathPrefix
	this.ZkBaseInit(conns, connTimeoutS, this)
	return nil
}

func (this *CWithcacheZkAdapter) SetCallback(callback ICallback, userData interface{}) {
	this.m_callback = callback
	this.m_callbackUserData = userData
}

func (this *CWithcacheZkAdapter) AddConnProperty(conn *proto.CConnectProperty) error {
	return this.AddConnProperty(conn)
}

func (this *CWithcacheZkAdapter) UpdateConnProperty(conn *proto.CConnectProperty) error {
	return this.UpdateConnProperty(conn)
}

func (this *CWithcacheZkAdapter) DeleteConnProperty(serviceId *string) error {
	return this.DeleteConnProperty(serviceId)
}

func (this *CWithcacheZkAdapter) AfterConnect(conn *zk.Conn) error {
	this.m_conn = conn
	return nil
}

func (this *CWithcacheZkAdapter) EventCallback(event zk.Event) {
	if event.Type == zk.EventNodeCreated || event.Type == zk.EventNodeDeleted || event.Type == zk.EventNodeDataChanged {
		fmt.Println("[INFO] node change")
	}
}

func (this *CWithcacheZkAdapter) GetMasterNode(serverName string) (*proto.CNodeData, error) {
	path := this.JoinPathPrefix(&this.m_pathPrefix, &serverName)
	childrens, _, err := this.m_conn.Children(*path)
	if err != nil {
		fmt.Println("[ERROR] get children error")
		return nil, err
	}
	isFind := false
	for _, child := range childrens {
		_, _, nodeName, err := this.SplitePath(child)
		if err != nil {
			fmt.Println("[WARNING] splitepath is error, path: ", child)
			continue
		}
		if *nodeName == proto.MasterNode {
			isFind = true
			break
		}
	}
	if isFind == false {
		return nil, errors.New("is not exist")
	}
	return nil, nil
}

func (this *CWithcacheZkAdapter) RoundRobin(serverName string) (*proto.CNodeData, error) {
	return nil, nil
}

func (this *CWithcacheZkAdapter) WeightRoundRobin(serverName string) (*proto.CNodeData, error) {
	return nil, nil
}

func (this *CWithcacheZkAdapter) Random(serverName string) (*proto.CNodeData, error) {
	return nil, nil
}

func (this *CWithcacheZkAdapter) WeightRandom(serverName string) (*proto.CNodeData, error) {
	return nil, nil
}

func (this *CWithcacheZkAdapter) IpHash(serverName string) (*proto.CNodeData, error) {
	return nil, nil
}

func (this *CWithcacheZkAdapter) UrlHash(serverName string) (*proto.CNodeData, error) {
	return nil, nil
}

func (this *CWithcacheZkAdapter) LeastConnections(serverName string) (*proto.CNodeData, error) {
	return nil, nil
}
