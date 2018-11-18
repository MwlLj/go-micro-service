package service_discovery_nocache

import (
	proto "../common_proto"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
	"time"
)

var _ = fmt.Println
var _ = errors.New("")

type CZkAdapter struct {
	m_zkCommon         proto.CZkCommon
	m_pathPrefix       string
	m_serverName       string
	m_serverUniqueCode string
	m_nodePayload      string
	m_connTimeout      int
	m_isConnected      bool
	m_callback         func(zk.Event)
	m_conn             *zk.Conn
}

type CNodeJson struct {
	ServerUniqueCode string `json:"serveruniquecode"`
	NodePayload      string `json:"nodepayload"`
}

func (this *CZkAdapter) init(conns *[]proto.CConnectProperty, serverName string, serverUniqueCode string, payload string, connTimeout int, pathPrefix string) error {
	this.m_serverName = serverName
	this.m_serverUniqueCode = serverUniqueCode
	this.m_nodePayload = payload
	this.m_connTimeout = connTimeout
	this.m_pathPrefix = pathPrefix
	this.m_isConnected = false
	for _, conn := range *conns {
		this.m_zkCommon.AddConnProperty(&conn)
	}
	this.m_callback = func(event zk.Event) {
		if event.State == zk.StateDisconnected {
			this.m_isConnected = false
		}
	}
	go func() {
		for {
			if this.m_isConnected == false {
				if this.m_conn != nil {
					this.m_conn.Close()
				}
				err := this.connect()
				if err == nil {
					this.m_isConnected = true
				}
			} else {
				time.Sleep(3 * time.Second)
			}
		}
	}()
	return nil
}

func (this *CZkAdapter) AddConnProperty(conn *proto.CConnectProperty) error {
	return this.m_zkCommon.AddConnProperty(conn)
}

func (this *CZkAdapter) UpdateConnProperty(conn *proto.CConnectProperty) error {
	return this.m_zkCommon.UpdateConnProperty(conn)
}

func (this *CZkAdapter) DeleteConnProperty(serviceId *string) error {
	return this.m_zkCommon.DeleteConnProperty(serviceId)
}

func (this *CZkAdapter) connect() error {
	option := zk.WithEventCallback(this.m_callback)
	hosts := this.m_zkCommon.ToHosts()
	var connChan <-chan zk.Event
	var err error = nil
	this.m_conn, connChan, err = zk.Connect(*hosts, time.Second, option)
	// this.m_conn, _, err = zk.Connect(*hosts, time.Second, option)
	if err != nil {
		fmt.Println("[ERROR] connect zookeeper server error")
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
	return this.createMasterAndNormalNode()
}

func (this *CZkAdapter) SetServerUniqueCode(uniqueCode string) {
	this.m_serverUniqueCode = uniqueCode
}

func (this *CZkAdapter) SetPayload(payload string) {
	this.m_nodePayload = payload
}

func (this *CZkAdapter) GetMasterPayload() (*string, error) {
	return nil, nil
}

func (this *CZkAdapter) joinNodeData() (*string, error) {
	nodeJson := CNodeJson{
		ServerUniqueCode: this.m_serverUniqueCode,
		NodePayload:      this.m_nodePayload}
	b, err := json.Marshal(&nodeJson)
	if err != nil {
		return nil, err
	}
	data := string(b)
	return &data, nil
}

func (this *CZkAdapter) createMasterAndNormalNode() error {
	data, err := this.joinNodeData()
	if err != nil {
		fmt.Println("[ERROR] join Node Data error")
		return err
	}
	pathPrefix := this.m_zkCommon.JoinPathPrefix(&this.m_pathPrefix, &this.m_serverName)
	err = this.createParents(*pathPrefix)
	if err != nil {
		fmt.Println("[ERROR] create parents error")
		return err
	}
	err = this.createMasterNode(data, pathPrefix)
	if err != nil {
		// master create error -> create normal node
		fmt.Println("[INFO] master node create error -> create normal node")
		err = this.createNormalNode(data, pathPrefix)
		if err == nil {
			fmt.Println("[SUCCESS] Identify: normal node")
		}
	} else {
		// master create success
		fmt.Println("[SUCCESS] Identify: master node")
	}
	return err
}

func (this *CZkAdapter) createParents(root string) error {
	this._createParents(root)
	this.createRootNode(root)
	return nil
}

func (this *CZkAdapter) _createParents(root string) error {
	isRoot, parent := this.m_zkCommon.GetParentNode(root)
	if isRoot == true {
		return nil
	} else {
		this._createParents(parent)
		this.createRootNode(parent)
	}
	return nil
}

func (this *CZkAdapter) createRootNode(path string) error {
	// node := strings.Join([]string{"/", path}, "")
	if path == "" {
		return nil
	}
	isExist, _, err := this.m_conn.Exists(path)
	if err != nil {
		fmt.Println("[ERROR] judge path is exist error", path)
		return err
	}
	if isExist {
		fmt.Println("[INFO] path already exist: ", path)
		return nil
	}
	fmt.Println("create path: ", path, isExist)
	_, err = this.m_conn.Create(path, nil, 0, zk.WorldACL(zk.PermAll))
	if err == nil {
		fmt.Println("create path success: ", path)
	}
	return err
}

func (this *CZkAdapter) createNode(path string, flag int32, payload *string) (afterCreateNode string, e error) {
	// node := strings.Join([]string{"/", path}, "")
	var data string = ""
	if payload != nil {
		data = *payload
	}
	return this.m_conn.Create(path, []byte(data), flag, zk.WorldACL(zk.PermAll))
}

func (this *CZkAdapter) checkNodeDelete(selfNode string, ech <-chan zk.Event) {
	for {
		select {
		case event := <-ech:
			fmt.Println("[INFO] check master is deleted")
			path := event.Path
			li := strings.Split(path, "/")
			length := len(li)
			if event.Type == zk.EventNodeDeleted && li[length-1] == proto.MasterNode {
				fmt.Println("[INFO] master node is deleted")
				data, err := this.joinNodeData()
				if err != nil {
					fmt.Println("[ERROR] join Node Data error")
					return
				}
				pathPrefix := this.m_zkCommon.JoinPathPrefix(&this.m_pathPrefix, &this.m_serverName)
				err = this.createParents(*pathPrefix)
				if err != nil {
					fmt.Println("[ERROR] create parents error")
					return
				}
				err = this.createMasterNode(data, pathPrefix)
				if err == nil {
					// master create success -> delete self
					fmt.Println("[INFO] master node create success -> delete self node")
					err = this.m_conn.Delete(selfNode, 0)
					if err != nil {
						fmt.Println("[ERROR] delete self node error: ", err)
						return
					} else {
						fmt.Println("[INFO] delete node, path: ", selfNode)
					}
					fmt.Println("[SUCCESS] Identify: master node")
				} else {
					this.listenMasterNode(&selfNode)
				}
				return
			}
		}
	}
}

func (this *CZkAdapter) createMasterNode(palyload *string, pathPrefix *string) error {
	masterPath := strings.Join([]string{*pathPrefix, proto.MasterNode}, "/")
	afterCreateNode, err := this.createNode(masterPath, zk.FlagEphemeral, palyload)
	if err == nil {
		fmt.Println("[INFO] create node path: ", afterCreateNode)
	}
	return err
}

func (this *CZkAdapter) createNormalNode(palyload *string, pathPrefix *string) error {
	normalPath := strings.Join([]string{*pathPrefix, proto.NormalNode}, "/")
	afterCreateNode, err := this.createNode(normalPath, 3, palyload)
	if err != nil {
		fmt.Println("[ERROR] create normal node error, path: ", afterCreateNode)
		return err
	} else {
		fmt.Println("[INFO] create normale node success, path: ", afterCreateNode)
	}
	err = this.listenMasterNode(&afterCreateNode)
	if err != nil {
		fmt.Println("[ERROR] listen master node error")
	}
	return err
}

func (this *CZkAdapter) listenMasterNode(nodePath *string) error {
	_, parent := this.m_zkCommon.GetParentNode(*nodePath)
	masterNode := strings.Join([]string{parent, proto.MasterNode}, "/")
	_, _, ech, err := this.m_conn.ExistsW(masterNode)
	if err != nil {
		fmt.Println("[ERROR] listen master node error")
		return err
	}
	go this.checkNodeDelete(*nodePath, ech)
	return nil
}
