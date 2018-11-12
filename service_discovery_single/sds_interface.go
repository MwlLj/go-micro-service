package service_discovery_single

import ()

type CServiceDiscoverySingle interface {
	Init(property *CInitProperty) error
	Connect(property *CConnectProperty) error
}

func New(property *CInitProperty) CServiceDiscoverySingle {
	if property.ServerMode == ServerModeZookeeper {
		return &CZkAdapter{}
	}
	return nil
}
