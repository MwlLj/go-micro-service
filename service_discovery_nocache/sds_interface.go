package service_discovery_nocache

import (
	proto "github.com/MwlLj/go-micro-service/common_proto"
)

var (
	ServerModeZookeeper = "zookeeper"
)

type CInitProperty struct {
	// zookeeper ...
	PathPrefix   string
	ServerMode   string
	ServerName   string
	NodeData     proto.CNodeData
	ConnTimeoutS int
	Conns        []proto.CConnectProperty
}

type CServiceDiscoveryNocache interface {
	SetNodeData(data *proto.CNodeData)
	GetMasterPayload() (*string, error)
	AddConnProperty(conn *proto.CConnectProperty) error
	UpdateConnProperty(conn *proto.CConnectProperty) error
	DeleteConnProperty(serviceId *string) error
	init(conns *[]proto.CConnectProperty, serverName string, nodeData *proto.CNodeData, connTimeoutS int, pathPrefix string) error
}

func New(property *CInitProperty) CServiceDiscoveryNocache {
	if property.ServerMode == ServerModeZookeeper {
		adapter := &CZkAdapter{}
		adapter.init(&property.Conns, property.ServerName, &property.NodeData, property.ConnTimeoutS, property.PathPrefix)
		return adapter
	}
	return nil
}
