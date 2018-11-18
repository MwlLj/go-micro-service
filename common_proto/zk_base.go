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
}

type CZkBase struct {
	m_zkCommon    CZkCommon
	m_connTimeout int
	m_isConnected bool
	m_callback    func(zk.Event)
	m_conn        *zk.Conn
}

func (this *CZkBase) GetZkConn() *zk.Conn {
	return this.m_conn
}

func (this *CZkBase) Init(connTimeout int) {
	this.m_connTimeout = connTimeout
	this.m_isConnected = false
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
}

func (this *CZkBase) AfterConnect(conn *zk.Conn) error {
	return nil
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
	return this.AfterConnect(this.m_conn)
}
