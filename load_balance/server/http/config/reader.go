package config

import (
	"bytes"
	"encoding/json"
	"github.com/MwlLj/gotools/configs"
	"log"
)

type CReader struct {
	configs.CConfigBase
}

func (this *CReader) Read(path string) (*CConfigInfo, error) {
	// default value
	configInfo := CConfigInfo{}
	// brokerRegisterInfo
	brokerRegisterInfo := CServiceDiscory{}
	var serviceDiscoryNet []CServiceDiscoryNet
	var net CServiceDiscoryNet
	net.ServerHost = "localhost"
	net.ServerPort = 2181
	net.ServiceId = "server_1"
	serviceDiscoryNet = append(serviceDiscoryNet, net)
	brokerRegisterInfo.Conns = serviceDiscoryNet
	var nodeData CBrokerNetInfo
	nodeData.Host = "localhost"
	nodeData.Port = 1883
	nodeData.ServerUniqueCode = "fcae5124-f37e-4987-a6d7-3782f90c6288"
	brokerRegisterInfo.NodeData = nodeData
	brokerRegisterInfo.PathPrefix = "micro-service"
	brokerRegisterInfo.ServerMode = "zookeeper"
	brokerRegisterInfo.ServerName = "mqtt-nginx"
	brokerRegisterInfo.Weight = 1
	configInfo.BrokerRegisterInfo = brokerRegisterInfo
	// serviceRegisterInfo
	var serviceLoadBalanceInfo CLoadBalanceInfo
	var serviceServiceDiscoryNet []CServiceDiscoryNet
	var serviceNet CServiceDiscoryNet
	serviceNet.ServerHost = "localhost"
	serviceNet.ServerPort = 2181
	serviceNet.ServiceId = "server_1"
	serviceServiceDiscoryNet = append(serviceServiceDiscoryNet, serviceNet)
	serviceLoadBalanceInfo.Conns = serviceServiceDiscoryNet
	serviceLoadBalanceInfo.PathPrefix = "face-service"
	serviceLoadBalanceInfo.ServerMode = "zookeeper_mqtt"
	serviceLoadBalanceInfo.NormalNodeAlgorithm = "roundrobin"
	configInfo.ServiceRegisterInfo = serviceLoadBalanceInfo
	// recvBrokerInfo
	recvBrokerInfo := CRecvBrokerInfo{}
	var recvNets []CBrokerNetInfo
	var recvNet CBrokerNetInfo
	recvNet.Host = "localhost"
	recvNet.Port = 1883
	recvNets = append(recvNets, recvNet)
	recvBrokerInfo.Nets = recvNets
	configInfo.RecvBrokerInfo = recvBrokerInfo
	// routerRuleInfo
	routerRuleInfo := CRouterRuleInfo{}
	var rules []CRouterRule
	var rule CRouterRule
	rule.Rule = ".*"
	rule.ServerName = "test"
	rule.IsMaster = false
	rules = append(rules, rule)
	routerRuleInfo.Rules = rules
	configInfo.RouterRuleInfo = routerRuleInfo
	// encoder
	b, err := json.Marshal(&configInfo)
	if err != nil {
		log.Fatalln("read config encoder error, err: ", err)
		return nil, err
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "\t")
	if err != nil {
		log.Fatalln("config file format change error, err: ", err)
		return nil, err
	}
	value, err := this.Load(path, out.String())
	if err != nil {
		log.Fatalln("read config error, err: ", err)
		return nil, err
	}
	// parse output
	outputConfig := CConfigInfo{}
	err = json.Unmarshal([]byte(value), &outputConfig)
	if err != nil {
		log.Fatalln("parer config error, err: ", err)
		return nil, err
	}
	return &outputConfig, nil
}
