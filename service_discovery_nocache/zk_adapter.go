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

type CZkAdapter struct {
	m_pathPrefix  string
	m_serverName  string
	m_connTimeout int
	m_connMap     sync.Map
}

type CConnInfo struct {
	host string
}

func (this *CZkAdapter) init(conns *[]CConnectProperty, serverName string, connTimeout int, pathPrefix string) error {
	this.m_serverName = serverName
	this.m_connTimeout = connTimeout
	this.m_pathPrefix = pathPrefix
	for _, conn := range *conns {
		this.AddConnProperty(&conn)
	}
	return nil
}

func (this *CZkAdapter) Connect() error {
	hosts := this.toHosts()
	conn, connChan, err := zk.Connect(*hosts, time.Second)
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
	isExist, err := this.createMasterNode(conn, "I'm is master")
	if err != nil {
		return err
	}
	if isExist {
		err = this.createNormalNode(conn, "I'm is normal")
	}
	return err
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

func (this *CZkAdapter) createParents(conn *zk.Conn, root string) error {
	this._createParents(conn, root)
	this.createNode(conn, root, 0, nil)
	return nil
}

func (this *CZkAdapter) _createParents(conn *zk.Conn, root string) error {
	isRoot, parent := this.getParentNode(root)
	if isRoot == true {
		return nil
	} else {
		this._createParents(conn, parent)
		this.createNode(conn, parent, 0, nil)
	}
	return nil
}

func (this *CZkAdapter) createNode(conn *zk.Conn, path string, flag int32, payload *string) (exist bool, e error) {
	node := strings.Join([]string{"/", path}, "")
	isExist, _, err := conn.Exists(node)
	if err != nil {
		return false, err
	}
	if isExist {
		return true, nil
	}
	fmt.Println("create node: ", node, isExist)
	_, err = conn.Create(node, []byte(*payload), flag, zk.WorldACL(zk.PermAll))
	if err == nil {
		fmt.Println("create node success: ", node)
	}
	return true, err
}

func (this *CZkAdapter) getParentNode(path string) (isRoot bool, parent string) {
	li := strings.Split(path, "/")
	length := len(li)
	if length == 1 {
		return true, li[0]
	}
	return false, strings.Join(li[:length-1], "/")
}

func (this *CZkAdapter) createMasterNode(conn *zk.Conn, payload string) (exist bool, e error) {
	var masterRoot string
	var masterPath string
	if this.m_pathPrefix != "" {
		masterRoot = strings.Join([]string{this.m_pathPrefix, this.m_serverName}, "/")
	} else {
		masterRoot = this.m_serverName
	}
	err := this.createParents(conn, masterRoot)
	if err != nil {
		return false, err
	}
	masterPath = strings.Join([]string{masterRoot, "master"}, "/")
	return this.createNode(conn, masterPath, zk.FlagEphemeral, &payload)
}

func (this *CZkAdapter) createNormalNode(conn *zk.Conn, payload string) error {
	var normalPath string
	if this.m_pathPrefix != "" {
		normalPath = strings.Join([]string{this.m_pathPrefix, this.m_serverName, "node"}, "/")
	} else {
		normalPath = strings.Join([]string{this.m_serverName, "node"}, "/")
	}
	_, err := this.createNode(conn, normalPath, 3, &payload)
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
