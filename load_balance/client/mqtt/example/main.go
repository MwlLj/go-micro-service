package main

import (
	bl "../../.."
	"../client"
	"../config"
	"fmt"
	"log"
)

var _ = log.Println

func main() {
	info := config.CConfigInfo{}
	loadBalanceInfo := config.CLoadBalanceInfo{}
	loadBalanceInfo.PathPrefix = "micro-service"
	loadBalanceInfo.ServerMode = bl.ServerModeZookeeper
	loadBalanceInfo.NormalNodeAlgorithm = bl.AlgorithmRoundRobin
	loadBalanceInfo.MqttLoadBalanceServerName = "mqtt-nginx"
	loadBalanceInfo.MqttServerName = "cgws"
	loadBalanceInfo.MqttServerVersion = "1.0"
	loadBalanceInfo.MqttServerRecvQos = 0
	var conns []config.CServiceDiscoryNet
	conn := config.CServiceDiscoryNet{
		ServerHost: "localhost",
		ServerPort: 2182,
		ServiceId:  "server_1",
	}
	conns = append(conns, conn)
	loadBalanceInfo.Conns = conns
	info.MqttLoadBalanceInfo = loadBalanceInfo
	cli := client.New(&info)
	mqttComm, urlMaker, err := cli.GetConnect()
	if err != nil {
		log.Fatalln(err)
		return
	}
	topic := "configs/serverinfo"
	fmt.Println(*urlMaker.Make(&topic))
	mqttComm.Post(*urlMaker.Make(&topic), string("hello"), 1, 30*60*1000*1000*1000)
}
