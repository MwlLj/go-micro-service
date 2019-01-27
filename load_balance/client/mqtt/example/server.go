package main

import (
	bl "../../.."
	s "../../../../service_discovery_nocache"
	"../client"
	"../config"
	"fmt"
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
	info := config.CConfigInfo{}
	loadBalanceInfo := config.CLoadBalanceInfo{}
	loadBalanceInfo.PathPrefix = "micro-service"
	loadBalanceInfo.ServerMode = bl.ServerModeZookeeper
	loadBalanceInfo.NormalNodeAlgorithm = bl.AlgorithmRoundRobin
	loadBalanceInfo.MqttLoadBalanceServerName = "mqtt-nginx"
	var conns []config.CServiceDiscoveryNet
	conn := config.CServiceDiscoveryNet{
		ServerHost: "localhost",
		ServerPort: 2182,
		ServiceId:  "server_1",
	}
	conns = append(conns, conn)
	loadBalanceInfo.Conns = conns
	info.MqttLoadBalanceInfo = loadBalanceInfo
	// service init
	uid, err := uuid.NewV4()
	if err != nil {
		fmt.Println(err)
		return
	}
	serverUniqueCode := uid.String()
	serviceInfo := config.CServiceInfo{}
	serviceInfo.ServerName = "configs"
	serviceInfo.ServerVersion = "1.0"
	serviceInfo.ServerRecvQos = 0
	var serviceServiceDiscoveryConns []config.CServiceDiscoveryNet
	serviceServiceDiscoveryConn := config.CServiceDiscoveryNet{
		ServerHost: "localhost",
		ServerPort: 2182,
		ServiceId:  "server_1",
	}
	serviceServiceDiscoveryConns = append(serviceServiceDiscoveryConns, serviceServiceDiscoveryConn)
	serviceInfo.ServiceDiscoveryConns = serviceServiceDiscoveryConns
	serviceInfo.PathPrefix = "taobao-service"
	serviceInfo.ServerMode = s.ServerModeZookeeper
	serviceInfo.ServerUniqueCode = serverUniqueCode
	serviceInfo.Weight = 1
	serviceInfo.BrokerHost = "localhost"
	serviceInfo.BrokerPort = 51883
	serviceInfo.BrokerUserName = ""
	serviceInfo.BrokerUserPwd = ""
	info.ServiceInfo = serviceInfo
	cli := client.New(&info)
	cli.Subscribe(mqtt_comm.POST, "configs/serverinfo", 0, &CAddServerInfoHandler{}, this)
	cli.StartRecver()
}

func main() {
	server := CServer{}
	server.Start()
}
