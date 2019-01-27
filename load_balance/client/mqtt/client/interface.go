package client

import (
	"../config"
	"../impl"
	"../url"
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
