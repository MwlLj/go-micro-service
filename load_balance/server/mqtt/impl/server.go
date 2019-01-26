package main

import (
	bl "../../../"
	proto "../../../../common_proto"
	sd "../../../../service_discovery_nocache"
	cfg "../config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

var _ = fmt.Println

type CServer struct {
	m_configReader cfg.CReader
}

func (this *CServer) Start(path string) error {
	configInfo, err := this.m_configReader.Read(path)
	if err != nil {
		log.Fatalln("mqtt load balance start failed, err: ", err)
		return err
	}
	// register mqtt load balance to service discovery
	brokerRegisterInfo := configInfo.BrokerRegisterInfo
	var brokerServiceDiscoveryConns []proto.CConnectProperty
	for _, item := range brokerRegisterInfo.Conns {
		brokerServiceDiscoveryConns = append(brokerServiceDiscoveryConns, proto.CConnectProperty{
			ServerHost: item.ServerHost,
			ServerPort: item.ServerPort,
			ServiceId:  item.ServerId,
		})
	}
	brokerServiceObj := sd.New(&sd.CInitProperty{
		PathPrefix: brokerRegisterInfo.PathPrefix,
		ServerMode: brokerRegisterInfo.ServerMode,
		ServerName: brokerRegisterInfo.ServerName,
		NodeData: proto.CNodeData{
			ServerIp:         brokerRegisterInfo.NodeData.Host,
			ServerPort:       brokerRegisterInfo.NodeData.Port,
			ServerUniqueCode: "",
			Weight:           brokerRegisterInfo.Weight,
		},
		Conns:        brokerServiceDiscoveryConns,
		ConnTimeoutS: 10})
	if brokerServiceObj == nil {
		log.Fatalln("mqtt load balance register service discovery error")
		return errors.New("register to service discovery error")
	}
	// register service to service discovery
	serviceRegisterInfo := configInfo.ServiceRegisterInfo
	var ServiceServiceDiscoveryConns []proto.CConnectProperty
	for _, item := range serviceRegisterInfo {
		ServiceServiceDiscoveryConns = append(ServiceServiceDiscoveryConns, proto.CConnectProperty{
			ServerHost: item.ServiceHost,
			ServerPort: item.ServicePort,
			ServiceId:  item.ServiceId,
		})
	}
	bls, connChan := bl.New(
		serviceRegisterInfo.ServerMode,
		&ServiceServiceDiscoveryConns,
		serviceRegisterInfo.PathPrefix,
		10,
	)
	select {
	case <-connChan:
		break
	}
	err = bls.SetNormalNodeAlgorithm(serviceRegisterInfo.NormalNodeAlgorithm)
	if err != nil {
		log.Fatalln("NormalNodeAlgorithm is not support, err: ", err)
		return err
	}
	// add recv broker info
	recvBrokerInfo := configInfo.RecvBrokerInfo
	for _, item := range recvBrokerInfo.Nets {
		bls.AddRecvNetInfo(nil, &bl.CNetInfo{
			Host:     item.Host,
			Port:     item.Port,
			UserName: item.UserName,
			UserPwd:  item.UserPwd,
		})
	}
	// add router rules
	routerRuleInfo := configInfo.RouterRuleInfo
	var blConfigInfo bl.CConfigInfo
	var rules map[string]*bl.CRuleInfo = make(map[string]*bl.CRuleInfo)
	for _, item := range routerRuleInfo.Rules {
		rule := bl.CConfigInfo
		rule.ObjServerName = item.ServerName
		rule.IsMater = item.IsMaster
		rules[item.Rule] = &rule
	}
	blConfigInfo.Rules = rules
	bls.SetConfigInfo(&blConfigInfo)
	// run
	err = bls.Run(bl.CNetInfo{
		Host: brokerRegisterInfo.NodeData.Host,
		Port: brokerRegisterInfo.NodeData.Port,
	})
	if err != nil {
		log.Fatalln("run mqtt load balance error, err: ", err)
		return err
	}
	return nil
}
