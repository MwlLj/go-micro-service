package main

import (
	"fmt"
	"github.com/MwlLj/mqtt_comm"
)

var _ = fmt.Println

type CServer struct {
	m_mqttComm mqtt_comm.CMqttComm
}

type CAddServerInfoHandler struct {
}

func (this *CAddServerInfoHandler) Handle(topic *string, action *string, request *string, qos int, mc mqtt_comm.CMqttComm, user interface{}) (*string, error) {
	return nil, nil
}

func (this *CServer) Start() {
	this.m_mqttComm = mqtt_comm.NewMqttComm("configs", "1.0", 0)
	this.m_mqttComm.SetMessageBus("localhost", 51883, "", "")
	this.m_mqttComm.Subscribe(mqtt_comm.POST, "1.0/configs/serverinfo", 0, &CAddServerInfoHandler{}, this)
	this.m_mqttComm.Connect(true)
}

func main() {
	server := CServer{}
	server.Start()
}
