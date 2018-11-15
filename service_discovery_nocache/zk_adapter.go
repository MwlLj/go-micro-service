package service_discovery_nocache

import (
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"strconv"
	"strings"
	"sync"
	"time"
)

var _ = fmt.Println

var (
	master string = "master"
)

type CZkAdapter struct {
	m_pathPrefix       string
	m_serverName       string
	m_serverUniqueCode string
	m_nodePayload      string
	m_connTimeout      int
	m_connMap          sync.Map
	m_conn             *zk.Conn
}

type CConnInfo struct {
	host string
}

type CNodeJson struct {
	ServerUniqueCode string `json:"serveruniquecode"`
	NodePayload      string `json:"nodepayload"`
}

func (this *CZkAdapter) init(conns *[]CConnectProperty, serverName string, serverUniqueCode string, payload string, connTimeout int, pathPrefix string) error {
	this.m_serverName = serverName
	this.m_serverUniqueCode = serverUniqueCode
	this.m_nodePayload = payload
	this.m_connTimeout = connTimeout
	this.m_pathPrefix = pathPrefix
	for _, conn := range *conns {
		this.AddConnProperty(&conn)
	}
	return nil
}

func (this *CZkAdapter) Connect() error {
	hosts := this.toHosts()
	var connChan <-chan zk.Event
	var err error = nil
	this.m_conn, connChan, err = zk.Connect(*hosts, time.Second)
	if err != nil {
		fmt.Println("connect zookeeper server error")
		return err
	}
	t := time.After(time.Second * time.Duration(this.m_connTimeout))
end:
	for {
		select {
		case event := <-connChan:
			if event.State == zk.StateConnected {
				fmt.Println("connect success")
				break end
			}
		case <-t:
			return errors.New("[Error] connect timeout")
		}
	}
	// connect success
	return this.createMasterAndNormalNode()
}

func (this *CZkAdapter) AddConnProperty(conn *CConnectProperty) error {
	info := CConnInfo{}
	info.host = this.joinHost(conn.ServerHost, conn.ServerPort)
	this.m_connMap.Store(conn.ServiceId, info)
	return nil
}

func (this *CZkAdapter) UpdateConnProperty(conn *CConnectProperty) error {
	info := CConnInfo{}
	info.host = this.joinHost(conn.ServerHost, conn.ServerPort)
	this.m_connMap.Store(conn.ServiceId, info)
	return nil
}

func (this *CZkAdapter) DeleteConnProperty(serviceId *string) error {
	this.m_connMap.Delete(serviceId)
	return nil
}

func (this *CZkAdapter) SetPayload(payload string) {
	this.m_nodePayload = payload
}

func (this *CZkAdapter) GetMasterPayload() (*string, error) {
	return nil, nil
}

func (this *CZkAdapter) createMasterAndNormalNode() error {
	isExist, err := this.createMasterNode("I'm is master")
	if isExist || err != nil {
		fmt.Println("master is not exist")
		err = this.createNormalNode("I'm is normal")
	}
	return err
}

func (this *CZkAdapter) createParents(root string) error {
	this._createParents(root)
	this.createNode(root, 0, nil, false)
	return nil
}

func (this *CZkAdapter) _createParents(root string) error {
	isRoot, parent := this.getParentNode(root)
	if isRoot == true {
		return nil
	} else {
		this._createParents(parent)
		this.createNode(parent, 0, nil, false)
	}
	return nil
}

func (this *CZkAdapter) createNode(path string, flag int32, payload *string, isListen bool) (exist bool, e error) {
	// return -> exist :  before create node is exist
	node := strings.Join([]string{"/", path}, "")
	isExist, _, err := this.m_conn.Exists(node)
	if err != nil {
		fmt.Println("judge node is exist error", node)
		return false, err
	}
	if isExist {
		fmt.Println("node already exist: ", node)
		return true, nil
	}
	fmt.Println("create node: ", node, isExist)
	afterCreateNode, err := this.m_conn.Create(node, []byte(*payload), flag, zk.WorldACL(zk.PermAll))
	if err == nil {
		fmt.Println("create node success: ", node)
	}
	// listen
	if isListen {
		_, parent := this.getParentNode(afterCreateNode)
		masterNode := strings.Join([]string{parent, master}, "/")
		fmt.Println("listen node: ", masterNode)
		_, _, ech, err := this.m_conn.ExistsW(masterNode)
		if err != nil {
			fmt.Println("listen error")
			return false, err
		}
		go this.checkNodeDelete(afterCreateNode, ech)
	}
	return false, err
}

func (this *CZkAdapter) checkNodeDelete(selfNode string, ech <-chan zk.Event) {
	for {
		select {
		case event := <-ech:
			path := event.Path
			li := strings.Split(path, "/")
			length := len(li)
			fmt.Println(path)
			if event.Type == zk.EventNodeDeleted && li[length-1] == master {
				fmt.Println("master deleted")
				err := this.m_conn.Delete(selfNode, 0)
				if err != nil {
					fmt.Println("delete self node error: ", err)
					return
				}
				err = this.createMasterAndNormalNode()
				if err != nil {
					fmt.Println("checkNodeDelete error: createMasterAndNormalNode error: ", err)
				}
				return
			}
		}
	}
}

func (this *CZkAdapter) getParentNode(path string) (isRoot bool, parent string) {
	li := strings.Split(path, "/")
	length := len(li)
	if length == 1 {
		return true, li[0]
	}
	return false, strings.Join(li[:length-1], "/")
}

func (this *CZkAdapter) createMasterNode(payload string) (exist bool, e error) {
	var masterRoot string
	var masterPath string
	if this.m_pathPrefix != "" {
		masterRoot = strings.Join([]string{this.m_pathPrefix, this.m_serverName}, "/")
	} else {
		masterRoot = this.m_serverName
	}
	err := this.createParents(masterRoot)
	if err != nil {
		return false, err
	}
	masterPath = strings.Join([]string{masterRoot, master}, "/")
	isExist, err := this.createNode(masterPath, zk.FlagEphemeral, &payload, false)
	if err == nil {
		fmt.Println("[SUCCESS] Identify: master node")
	}
	return isExist, err
}

func (this *CZkAdapter) createNormalNode(payload string) error {
	var normalPath string
	if this.m_pathPrefix != "" {
		normalPath = strings.Join([]string{this.m_pathPrefix, this.m_serverName, "node"}, "/")
	} else {
		normalPath = strings.Join([]string{this.m_serverName, "node"}, "/")
	}
	_, err := this.createNode(normalPath, 3, &payload, true)
	if err == nil {
		fmt.Println("[SUCCESS] Identify: normal node")
	}
	return err
}

func (*CZkAdapter) joinHost(ip string, port int) string {
	return strings.Join([]string{ip, strconv.FormatInt(int64(port), 10)}, ":")
}

func (this *CZkAdapter) toHosts() *[]string {
	var hosts []string
	f := func(k, v interface{}) bool {
		info := v.(CConnInfo)
		hosts = append(hosts, info.host)
		return true
	}
	this.m_connMap.Range(f)
	return &hosts
}
