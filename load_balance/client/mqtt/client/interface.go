package client

import (
	"github.com/MwlLj/go-micro-service/load_balance/client/mqtt/config"
	"github.com/MwlLj/go-micro-service/load_balance/client/mqtt/impl"
	"github.com/MwlLj/go-micro-service/load_balance/client/mqtt/url"
	"github.com/MwlLj/mqtt_comm"
)

type IClient interface {
	GetConnect() (mqtt_comm.CMqttComm, url.IUrlMaker, error)
	Subscribe(action string, topic string, qos int, handler mqtt_comm.CHandler, user interface{}) error
	StartRecver()
}

func New(info *config.CConfigInfo) IClient {
	client := impl.CClient{}
	client.Init(info)
	return &client
}
