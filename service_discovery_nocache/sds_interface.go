package service_discovery_nocache

import ()

var (
	ServerModeZookeeper = "zookeeper"
)

type CInitProperty struct {
	// zookeeper ...
	PathPrefix   string
	ServerMode   string
	ServerName   string
	ConnTimeoutS int
	Conns        []CConnectProperty
}

type CConnectProperty struct {
	ServerHost string
	ServerPort int
	ServiceId  string
}

type CServiceDiscoveryNocache interface {
	Connect() error
	AddConnProperty(conn *CConnectProperty) error
	UpdateConnProperty(conn *CConnectProperty) error
	DeleteConnProperty(serviceId *string) error
	init(conns *[]CConnectProperty, serverName string, connTimeout int, pathPrefix string) error
}

func New(property *CInitProperty) CServiceDiscoveryNocache {
	if property.ServerMode == ServerModeZookeeper {
		adapter := &CZkAdapter{}
		adapter.init(&property.Conns, property.ServerName, property.ConnTimeoutS, property.PathPrefix)
		return adapter
	}
	return nil
}
