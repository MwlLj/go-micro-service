package service_discovery_single

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"strconv"
	"strings"
)

var _ = fmt.Println

type CZkAdapter struct {
	m_conn *zk.Conn
}

func (this *CZkAdapter) Init(property *CInitProperty) error {
	return nil
}

func (this *CZkAdapter) Connect(property *CConnectProperty) error {
	host := strings.Join([]string{property.ZkServerHost, strconv.FormatInt(property.ZkServerPort, 10)}, ":")
	this.m_conn, connChan, err = zk.Connect(host, property.ConnTimeoutS)
	if err != nil {
		fmt.Println("connect zookeeper server error")
		return err
	}
	return nil
}
