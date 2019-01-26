package client

import (
	"../config"
	"../impl"
	"github.com/MwlLj/mqtt_comm"
)

type IClient interface {
	GetConnect() mqtt_comm.CMqttComm
	JoinTopic(topic string) *string
}

func New(info *config.CConfigInfo) IClient {
	client := impl.CClient{}
	client.Init(info)
	return &client
}
