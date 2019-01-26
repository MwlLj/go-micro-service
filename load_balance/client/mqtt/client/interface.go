package client

import (
	"../config"
	"../impl"
	"github.com/MwlLj/mqtt_comm"
)

type IClient interface {
	GetConnect() (mqtt_comm.CMqttComm, error)
	JoinTopic(topic string, serverUniqueCode *string) *string
	Subscribe(action string, topic string, qos int, handler mqtt_comm.CHandler, user interface{}) error
}

func New(info *config.CConfigInfo) IClient {
	client := impl.CClient{}
	client.Init(info)
	return &client
}
