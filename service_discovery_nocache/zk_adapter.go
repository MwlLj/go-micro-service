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
	m_connTimeout int
	m_connMap     sync.Map
}

type CConnInfo struct {
	host string
}

func (this *CZkAdapter) init(conns *[]CConnectProperty, connTimeout int) error {
	this.m_connTimeout = connTimeout
	for _, conn := range *conns {
		this.AddConnProperty(&conn)
	}
	return nil
}

func (this *CZkAdapter) Connect() error {
	hosts := this.toHosts()
	conn, connChan, err := zk.Connect(*hosts, time.Duration(this.m_connTimeout)*time.Second)
	if err != nil {
		fmt.Println("connect zookeeper server error")
		return err
	}
end:
	for {
		select {
		case event := <-connChan:
			if event.State == zk.StateConnected {
				fmt.Println("connect success")
				break end
			}
		case <-time.After(time.Second * 3):
			return errors.New("[Error] connect timeout")
		}
	}
	_ = conn
	return nil
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
