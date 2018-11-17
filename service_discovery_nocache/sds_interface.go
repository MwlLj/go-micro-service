package service_discovery_nocache

import (
	proto "../common_proto"
)

var (
	ServerModeZookeeper = "zookeeper"
)

type CInitProperty struct {
	// zookeeper ...
	PathPrefix       string
	ServerMode       string
	ServerName       string
	ServerUniqueCode string
	NodePayload      string
	ConnTimeoutS     int
	Conns            []proto.CConnectProperty
}

type CServiceDiscoveryNocache interface {
	Connect() error
	SetServerUniqueCode(uniqueCode string)
	SetPayload(payload string)
	GetMasterPayload() (*string, error)
	AddConnProperty(conn *proto.CConnectProperty) error
	UpdateConnProperty(conn *proto.CConnectProperty) error
	DeleteConnProperty(serviceId *string) error
	init(conns *[]proto.CConnectProperty, serverName string, serverUniqueCode string, payload string, connTimeout int, pathPrefix string) error
}

func New(property *CInitProperty) CServiceDiscoveryNocache {
	if property.ServerMode == ServerModeZookeeper {
		adapter := &CZkAdapter{}
		adapter.init(&property.Conns, property.ServerName, property.ServerUniqueCode, property.NodePayload, property.ConnTimeoutS, property.PathPrefix)
		return adapter
	}
	return nil
}
