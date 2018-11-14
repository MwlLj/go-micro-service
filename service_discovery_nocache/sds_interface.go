package service_discovery_nocache

import ()

var (
	ServerModeZookeeper = "zookeeper"
)

type CInitProperty struct {
	// zookeeper ...
	ServerMode   string
	ConnProperty CConnectProperty
}

type CNet struct {
	ServerHost   string
	ServerPort   int
	ConnTimeoutS int
	ServiceName  string
}

type CConnectProperty struct {
	Nets []CNet
}

type CServiceDiscoveryNocache interface {
	Init(property *CInitProperty) error
	Connect() error
}

func New(property *CInitProperty) CServiceDiscoveryNocache {
	if property.ServerMode == ServerModeZookeeper {
		adapter := &CZkAdapter{}
		adapter.Init(property)
		return adapter
	}
	return nil
}
