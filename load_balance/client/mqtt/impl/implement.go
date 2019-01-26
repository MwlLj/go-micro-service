package impl

import (
	bl "../../.."
	proto "../../../../common_proto"
	sd "../../../../service_discovery_nocache"
	"../config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MwlLj/mqtt_comm"
	"log"
	"strconv"
	"sync"
)

var _ = fmt.Println
var _ = bytes.Equal
var _ = json.Marshal

type CClinet struct {
	m_configInfo  *config.CConfigInfo
	m_balance     bl.ILoadBlance
	m_algorithm   bl.INormalNodeAlgorithm
	m_mqttCommMap sync.Map
}

func (this *CClinet) GetConnect() (mqtt_comm.CMqttComm, error) {
	nodeItem, err := this.m_balance.FindServerData(this.m_configInfo.MqttLoadBalanceServerName)
	if err != nil {
		log.Println("find server data from service discovery error, err: ", err)
		return nil, err
	}
	var host *string = nil
	if nodeItem.NormalNodes == nil || len(nodeItem.NormalNodes) == 0 {
		if nodeItem.MasterNode == nil {
			log.Println("normal node or master node both is not exist")
			return nil, errors.New("normal and master node is not exist")
		} else {
			host = this.joinHost(&nodeItem.MasterNode.ServerIp, nodeItem.MasterNode.ServerPort)
		}
	} else {
		if this.m_algorithm == nil {
			log.Println("normal node algorithm is nil")
			return nil, errors.New("normal node algorithm is nil")
		}
		data, err = this.m_algorithm.Get(this.m_configInfo.MqttLoadBalanceServerName, nil)
		if err != nil {
			log.Println("get normal node error, err: ", err)
			return nil, err
		}
		host = this.joinHost(&data.ServerIp, data.ServerPort)
	}
	var mqttComm mqtt_comm.CMqttComm
	value, ok := this.m_mqttCommMap.Load(*host)
	if !ok {
		// not exist
		mqttComm = mqtt_comm.NewMqttComm(this.m_configInfo.MqttServerName, this.m_configInfo.MqttServerVersion, this.m_configInfo.MqttServerRecvQos)
		this.m_mqttCommMap.Store(*host, mqttComm)
	} else {
		// exist
		mqttComm = value.(mqtt_comm.CMqttComm)
		/*
			if mqttComm.IsConnect == false {
				mqttComm = nil
				mqttComm = mqtt_comm.NewMqttComm(this.m_configInfo.MqttServerName, this.m_configInfo.MqttServerVersion, this.m_configInfo.MqttServerRecvQos)
			}
			this.m_mqttCommMap.Store(*host, mqttComm)
		*/
	}
	return mqttComm, nil
}

func (this *CClinet) JoinTopic(topic string) *string {
	var buffer bytes.Buffer
	buffer.WriteString(topic)
	bTopic := []byte(topic)
	if bTopic[len(bTopic)-1] != '/' {
		buffer.WriteString("/")
	}
	buffer.WriteString(*serverUniqueCode)
	top := buffer.String()
	return &top
}

func (this *CClinet) joinHost(host *string, port int) *string {
	var buffer bytes.Buffer
	buffer.WriteString(*host)
	buffer.WriteString(":")
	buffer.WriteString(strconv.FormatInt(int64(port), 10))
	h := buffer.String()
	return &h
}

func (this *CClient) Init(info *config.CConfigInfo) {
	if info == nil {
		log.Fatalln("client init error, info is nil")
		return
	}
	this.m_configInfo = info
	mqttLoadBalanceInfo := info.MqttLoadBalanceInfo
	var conns []proto.CConnectProperty
	for _, item := range mqttLoadBalanceInfo.Conns {
		conns = append(conns, proto.CConnectProperty{
			ServerHost: item.ServerHost,
			ServerPort: item.ServerPort,
			ServiceId:  item.ServiceId,
		})
	}
	var connChan <-chan bool
	this.m_balance, connChan = bl.New(
		serviceRegisterInfo.ServerMode,
		&conns,
		serviceRegisterInfo.PathPrefix,
		10,
	)
	select {
	case <-connChan:
		break
	}
	this.m_algorithm = bls.GetNormalNodeAlgorithm(mqttLoadBalanceInfo.NormalNodeAlgorithm)
	if this.m_algorithm == nil {
		log.Fatalln("get normalenode algorithm error")
		return
	}
}
