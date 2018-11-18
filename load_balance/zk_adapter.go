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

type CZkAdapter struct {
	proto.CZkBase
	m_conn *zk.Conn
}

func (this *CZkAdapter) init(conns *[]proto.CConnectProperty, connTimeoutS int) error {
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
	return nil
}

func (this *CZkAdapter) EventCallback(event zk.Event) {
}

func (this *CZkAdapter) GetMasterNode(serverName string) *string {
	return nil
}

func (this *CZkAdapter) RoundRobin(serverName string) *string {
	return nil
}

func (this *CZkAdapter) WeightRoundRobin(serverName string) *string {
	return nil
}

func (this *CZkAdapter) Random(serverName string) *string {
	return nil
}

func (this *CZkAdapter) WeightRandom(serverName string) *string {
	return nil
}

func (this *CZkAdapter) IpHash(serverName string) *string {
	return nil
}

func (this *CZkAdapter) UrlHash(serverName string) *string {
	return nil
}

func (this *CZkAdapter) LeastConnections(serverName string) *string {
	return nil
}
