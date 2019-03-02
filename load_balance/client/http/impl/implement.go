package impl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	proto "github.com/MwlLj/go-micro-service/common_proto"
	bl "github.com/MwlLj/go-micro-service/load_balance"
	"github.com/MwlLj/go-micro-service/load_balance/client/http/config"
	s "github.com/MwlLj/go-micro-service/service_discovery_nocache"
	"log"
)

var _ = fmt.Println
var _ = bytes.Equal
var _ = json.Marshal

type CClient struct {
	m_configInfo *config.CConfigInfo
	m_balance    bl.ILoadBlance
	m_algorithm  bl.INormalNodeAlgorithm
}

func (this *CClient) GetConnect() (*config.CLoadBalanceNetInfo, error) {
	nodeItem, err := this.m_balance.FindServerData(this.m_configInfo.HttpLoadBalanceInfo.HttpLoadBalanceServerName)
	if err != nil {
		log.Println("find server data from service discovery error, err: ", err)
		return nil, err
	}
	var loadBalanceNetInfo config.CLoadBalanceNetInfo
	if nodeItem.NormalNodes == nil || len(*nodeItem.NormalNodes) == 0 {
		if nodeItem.MasterNode == nil {
			log.Println("normal node or master node both is not exist")
			return nil, errors.New("normal and master node is not exist")
		} else {
			loadBalanceNetInfo.Host = nodeItem.MasterNode.ServerIp
			loadBalanceNetInfo.Port = nodeItem.MasterNode.ServerPort
			loadBalanceNetInfo.UserName = nodeItem.MasterNode.UserName
			loadBalanceNetInfo.UserPwd = nodeItem.MasterNode.UserPwd
			loadBalanceNetInfo.ServerUniqueCode = nodeItem.MasterNode.ServerUniqueCode
		}
	} else {
		if this.m_algorithm == nil {
			log.Println("normal node algorithm is nil")
			return nil, errors.New("normal node algorithm is nil")
		}
		data, err := this.m_algorithm.Get(this.m_configInfo.HttpLoadBalanceInfo.HttpLoadBalanceServerName, nil)
		if err != nil {
			log.Println("get normal node error, err: ", err)
			return nil, err
		}
		loadBalanceNetInfo.Host = data.ServerIp
		loadBalanceNetInfo.Port = data.ServerPort
		loadBalanceNetInfo.UserName = data.UserName
		loadBalanceNetInfo.UserPwd = data.UserPwd
		loadBalanceNetInfo.ServerUniqueCode = data.ServerUniqueCode
	}
	return &loadBalanceNetInfo, nil
}

func (this *CClient) RegisterService() error {
	var conns []proto.CConnectProperty
	for _, item := range this.m_configInfo.ServiceInfo.ServiceDiscoveryConns {
		conns = append(conns, proto.CConnectProperty{
			ServerHost: item.ServerHost,
			ServerPort: item.ServerPort,
			ServiceId:  item.ServiceId,
		})
	}
	sds := s.New(&s.CInitProperty{
		PathPrefix: this.m_configInfo.ServiceInfo.PathPrefix,
		ServerMode: this.m_configInfo.ServiceInfo.ServerMode,
		ServerName: this.m_configInfo.ServiceInfo.ServerName,
		NodeData: proto.CNodeData{
			ServerIp:         this.m_configInfo.ServiceInfo.Host,
			ServerPort:       this.m_configInfo.ServiceInfo.Port,
			ServerUniqueCode: this.m_configInfo.ServiceInfo.ServerUniqueCode,
			Weight:           this.m_configInfo.ServiceInfo.Weight,
		},
		Conns:        conns,
		ConnTimeoutS: 10})
	if sds == nil {
		log.Fatalln("connect service discovery error")
		return errors.New("connect service discovery error")
	}
	return nil
}

func (this *CClient) Init(info *config.CConfigInfo) {
	if info == nil {
		log.Fatalln("client init error, info is nil")
		return
	}
	this.m_configInfo = info
	loadBalanceInfo := info.HttpLoadBalanceInfo
	var conns []proto.CConnectProperty
	for _, item := range loadBalanceInfo.Conns {
		conns = append(conns, proto.CConnectProperty{
			ServerHost: item.ServerHost,
			ServerPort: item.ServerPort,
			ServiceId:  item.ServiceId,
		})
	}
	var connChan <-chan bool
	this.m_balance, connChan = bl.New(
		loadBalanceInfo.ServerMode,
		&conns,
		loadBalanceInfo.PathPrefix,
		10,
	)
	select {
	case <-connChan:
		break
	}
	this.m_algorithm = this.m_balance.GetNormalNodeAlgorithm(loadBalanceInfo.NormalNodeAlgorithm)
	if this.m_algorithm == nil {
		log.Fatalln("get normalenode algorithm error")
		return
	}
}
