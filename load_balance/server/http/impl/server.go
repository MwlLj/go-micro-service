package impl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	proto "github.com/MwlLj/go-micro-service/common_proto"
	bl "github.com/MwlLj/go-micro-service/load_balance"
	cfg "github.com/MwlLj/go-micro-service/load_balance/server/http/config"
	sd "github.com/MwlLj/go-micro-service/service_discovery_nocache"
	"log"
)

var _ = fmt.Println
var _ = bytes.Equal
var _ = json.Marshal

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
			ServiceId:  item.ServiceId,
		})
	}
	brokerServiceObj := sd.New(&sd.CInitProperty{
		PathPrefix: brokerRegisterInfo.PathPrefix,
		ServerMode: brokerRegisterInfo.ServerMode,
		ServerName: brokerRegisterInfo.ServerName,
		NodeData: proto.CNodeData{
			ServerIp:         brokerRegisterInfo.NodeData.Host,
			ServerPort:       brokerRegisterInfo.NodeData.Port,
			ServerUniqueCode: brokerRegisterInfo.NodeData.ServerUniqueCode,
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
	for _, item := range serviceRegisterInfo.Conns {
		ServiceServiceDiscoveryConns = append(ServiceServiceDiscoveryConns, proto.CConnectProperty{
			ServerHost: item.ServerHost,
			ServerPort: item.ServerPort,
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
		rule := bl.CRuleInfo{}
		rule.ObjServerName = item.ServerName
		rule.IsMaster = item.IsMaster
		rules[item.Rule] = &rule
	}
	blConfigInfo.Rules = rules
	bls.SetConfigInfo(&blConfigInfo)
	// run
	err = bls.Run(bl.CNetInfo{
		Host:       brokerRegisterInfo.NodeData.Host,
		Port:       brokerRegisterInfo.NodeData.Port,
		ExtraField: brokerRegisterInfo.NodeData.ServerUniqueCode,
	})
	if err != nil {
		log.Fatalln("run mqtt load balance error, err: ", err)
		return err
	}
	return nil
}
