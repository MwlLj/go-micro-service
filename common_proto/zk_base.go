package common_proto

import (
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

var _ = fmt.Println

type IZkBase interface {
	AfterConnect(conn *zk.Conn) error
	EventCallback(event zk.Event)
}

type CZkBase struct {
	m_baseInterface IZkBase
	m_zkCommon      CZkCommon
	m_connTimeout   int
	m_isConnected   bool
	m_callback      func(zk.Event)
	m_conn          *zk.Conn
}

func (this *CZkBase) ZkBaseInit(conns *[]CConnectProperty, connTimeout int, baseInterface IZkBase) {
	this.m_baseInterface = baseInterface
	for _, conn := range *conns {
		this.m_zkCommon.AddConnProperty(&conn)
	}
	this.m_connTimeout = connTimeout
	this.m_isConnected = false
	this.m_callback = func(event zk.Event) {
		fmt.Println("[INFO] watch event",
			"path:", event.Path,
			"state:", event.State,
			"type:", event.Type,
			"server:", event.Server)
		if event.State == zk.StateDisconnected {
			this.m_isConnected = false
		}
		if this.m_baseInterface != nil {
			this.m_baseInterface.EventCallback(event)
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
				} else {
					time.Sleep(1 * time.Second)
				}
			} else {
				time.Sleep(3 * time.Second)
			}
		}
	}()
}

func (this *CZkBase) connect() error {
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
	return this.m_baseInterface.AfterConnect(this.m_conn)
}

func (this *CZkBase) AddConnProperty(conn *CConnectProperty) error {
	return this.m_zkCommon.AddConnProperty(conn)
}

func (this *CZkBase) UpdateConnProperty(conn *CConnectProperty) error {
	return this.m_zkCommon.UpdateConnProperty(conn)
}

func (this *CZkBase) DeleteConnProperty(serviceId *string) error {
	return this.m_zkCommon.DeleteConnProperty(serviceId)
}

func (this *CZkBase) JoinPathPrefix(prefix *string, serverName *string) *string {
	return this.m_zkCommon.JoinPathPrefix(prefix, serverName)
}

func (this *CZkBase) GetParentNode(path string) (isRoot bool, parent string) {
	return this.m_zkCommon.GetParentNode(path)
}

func (this *CZkBase) SplitePath(path string) (prefix, serverName, nodeName *string, e error) {
	return this.m_zkCommon.SplitePath(path)
}
