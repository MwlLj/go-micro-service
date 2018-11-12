package service_discovery_single

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
)

var _ = fmt.Println

type CZkAdapter struct {
}

func (this *CZkAdapter) Init(property *CInitProperty) error {
	return nil
}

func (this *CZkAdapter) Connect(property *CConnectProperty) error {
	return nil
}
