package config

import (
	"bytes"
	"encoding/json"
	"github.com/MwlLj/gotools/configs"
	"github.com/satori/go.uuid"
	"log"
)

type CReader struct {
	configs.CConfigBase
}

func (this *CReader) Read(path string) (*CConfigInfo, error) {
	// gen loadbalance server uuid
	uid, err := uuid.NewV4()
	if err != nil {
		log.Fatalln("general server uuid error")
		return nil, err
	}
	serverUuid := uid.String()
	// default value
	configInfo := CConfigInfo{}
	// brokerRegisterInfo
	loadBalanceRegisterInfo := CServiceDiscory{}
	var serviceDiscoryNet []CServiceDiscoryNet
	var net CServiceDiscoryNet
	net.ServerHost = "localhost"
	net.ServerPort = 2181
	net.ServiceId = "server_1"
	serviceDiscoryNet = append(serviceDiscoryNet, net)
	loadBalanceRegisterInfo.Conns = serviceDiscoryNet
	var nodeData CNodeNetInfo
	nodeData.Host = "localhost"
	nodeData.Port = 30000
	nodeData.ServerUniqueCode = serverUuid
	loadBalanceRegisterInfo.NodeData = nodeData
	loadBalanceRegisterInfo.PathPrefix = "micro-service"
	loadBalanceRegisterInfo.ServerMode = "zookeeper"
	loadBalanceRegisterInfo.ServerName = "http-nginx"
	loadBalanceRegisterInfo.Weight = 1
	configInfo.LoadBalanceInfo = loadBalanceRegisterInfo
	// serviceRegisterInfo
	var serviceLoadBalanceInfo CLoadBalanceInfo
	var serviceServiceDiscoryNet []CServiceDiscoryNet
	var serviceNet CServiceDiscoryNet
	serviceNet.ServerHost = "localhost"
	serviceNet.ServerPort = 2181
	serviceNet.ServiceId = "server_1"
	serviceServiceDiscoryNet = append(serviceServiceDiscoryNet, serviceNet)
	serviceLoadBalanceInfo.Conns = serviceServiceDiscoryNet
	serviceLoadBalanceInfo.PathPrefix = "gateway-service"
	serviceLoadBalanceInfo.ServerMode = "zookeeper_http"
	serviceLoadBalanceInfo.NormalNodeAlgorithm = "roundrobin"
	configInfo.ServiceRegisterInfo = serviceLoadBalanceInfo
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
