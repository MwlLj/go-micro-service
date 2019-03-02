package client

import (
	"github.com/MwlLj/go-micro-service/load_balance/client/http/config"
	"github.com/MwlLj/go-micro-service/load_balance/client/http/impl"
)

type IClient interface {
	/*
		@bref get http load-balance server info
		@desc select load-balance
	*/
	GetConnect() (*config.CLoadBalanceNetInfo, error)
	/*
		@bref register service to register center
	*/
	RegisterService() error
}

func New(info *config.CConfigInfo) IClient {
	client := impl.CClient{}
	client.Init(info)
	return &client
}
