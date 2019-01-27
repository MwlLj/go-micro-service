package main

import (
	"fmt"
	proto "github.com/MwlLj/go-micro-service/common_proto"
	bl "github.com/MwlLj/go-micro-service/load_balance"
	s "github.com/MwlLj/go-micro-service/service_discovery_nocache"
	"github.com/MwlLj/mqtt_comm"
	"github.com/satori/go.uuid"
)

var _ = fmt.Println

type CServer struct {
	m_mqttComm mqtt_comm.CMqttComm
}

type CAddServerInfoHandler struct {
}

func (this *CAddServerInfoHandler) Handle(topic *string, action *string, request *string, qos int, mc mqtt_comm.CMqttComm, user interface{}) (*string, error) {
	fmt.Println("configs recv add serverinfo request ...")
	return nil, nil
}

func (this *CServer) Start() {
	var conns []proto.CConnectProperty
	conns = append(conns, proto.CConnectProperty{
		ServerHost: "127.0.0.1",
		ServerPort: 2182,
		ServiceId:  "server_1",
	})
	uid, err := uuid.NewV4()
	if err != nil {
		fmt.Println(err)
		return
	}
	serverUniqueCode := uid.String()
	sds := s.New(&s.CInitProperty{
		PathPrefix: "taobao-service",
		ServerMode: s.ServerModeZookeeper,
		ServerName: "configs",
		NodeData: proto.CNodeData{
			ServerIp:         "127.0.0.1",
			ServerPort:       50000,
			ServerUniqueCode: serverUniqueCode,
			Weight:           1,
		},
		Conns:        conns,
		ConnTimeoutS: 10})
	var _ = sds
	bls, connChan := bl.New(bl.ServerModeZookeeperMqtt, &conns, "micro-service", 10)
	select {
	case <-connChan:
		break
	}
	this.m_mqttComm = mqtt_comm.NewMqttComm("configs", "1.0", 0)
	this.m_mqttComm.SetMessageBus("localhost", 51883, "", "")
	topic := "configs/serverinfo"
	afterJoin := bls.TopicJoin(&topic, &serverUniqueCode)
	this.m_mqttComm.Subscribe(mqtt_comm.POST, *afterJoin, 0, &CAddServerInfoHandler{}, this)
	this.m_mqttComm.Connect(true)
}

func main() {
	server := CServer{}
	server.Start()
}
