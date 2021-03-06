package service_discovery_nocache

import (
	"encoding/json"
	"errors"
	"fmt"
	proto "github.com/MwlLj/go-micro-service/common_proto"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
)

var _ = fmt.Println
var _ = errors.New("")

type CZkAdapter struct {
	proto.CZkBase
	m_pathPrefix   string
	m_serverName   string
	m_isDisconnect bool
	m_nodeData     proto.CNodeData
	m_conn         *zk.Conn
}

func (this *CZkAdapter) init(conns *[]proto.CConnectProperty, serverName string, nodeData *proto.CNodeData, connTimeoutS int, pathPrefix string) error {
	this.m_serverName = serverName
	this.m_nodeData = *nodeData
	this.m_pathPrefix = pathPrefix
	this.m_isDisconnect = true
	this.ZkBaseInit(conns, connTimeoutS, this)
	return nil
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
	path, err := this.createMasterAndNormalNode()
	if err == nil {
		this.m_isDisconnect = false
	}
	go func() {
		for {
			_, _, ech, err := this.m_conn.ChildrenW(*path)
			if err != nil {
				continue
			}
			event := <-ech
			fmt.Println("[INFO] watch master/normal node children",
				"path:", event.Path,
				"state:", event.State,
				"type:", event.Type,
				"server:", event.Server)
		}
	}()
	return err
}

func (this *CZkAdapter) EventCallback(event zk.Event) {
	if event.State == zk.StateDisconnected {
		this.m_isDisconnect = true
	}
}

func (this *CZkAdapter) SetNodeData(data *proto.CNodeData) {
	this.m_nodeData = *data
}

func (this *CZkAdapter) GetMasterPayload() (*string, error) {
	return nil, nil
}

func (this *CZkAdapter) joinNodeData() (*string, error) {
	b, err := json.Marshal(&this.m_nodeData)
	if err != nil {
		return nil, err
	}
	d := string(b)
	return &d, nil
}

func (this *CZkAdapter) createMasterAndNormalNode() (*string, error) {
	data, err := this.joinNodeData()
	if err != nil {
		fmt.Println("[ERROR] join Node Data error")
		return nil, err
	}
	pathPrefix := this.JoinPathPrefix(&this.m_pathPrefix, &this.m_serverName)
	err = this.createParents(*pathPrefix)
	if err != nil {
		fmt.Println("[ERROR] create parents error")
		return nil, err
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
	return pathPrefix, err
}

func (this *CZkAdapter) createParents(root string) error {
	this._createParents(root)
	this.createRootNode(root)
	return nil
}

func (this *CZkAdapter) _createParents(root string) error {
	isRoot, parent := this.GetParentNode(root)
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
			if this.m_isDisconnect == true {
				return
			}
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
				pathPrefix := this.JoinPathPrefix(&this.m_pathPrefix, &this.m_serverName)
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
	_, parent := this.GetParentNode(*nodePath)
	masterNode := strings.Join([]string{parent, proto.MasterNode}, "/")
	_, _, ech, err := this.m_conn.ExistsW(masterNode)
	if err != nil {
		fmt.Println("[ERROR] listen master node error")
		return err
	}
	go this.checkNodeDelete(*nodePath, ech)
	return nil
}
