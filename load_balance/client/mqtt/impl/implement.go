package impl

import (
	bl "../../.."
	proto "../../../../common_proto"
	"../config"
	"../url"
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

type CSubscribeInfo struct {
	action   string
	topic    string
	qos      int
	handler  mqtt_comm.CHandler
	userData interface{}
}

type CClient struct {
	m_configInfo    *config.CConfigInfo
	m_balance       bl.ILoadBlance
	m_algorithm     bl.INormalNodeAlgorithm
	m_mqttCommMap   sync.Map
	m_subscribeInfo []*CSubscribeInfo
}

func (this *CClient) GetConnect() (mqtt_comm.CMqttComm, url.IUrlMaker, error) {
	nodeItem, err := this.m_balance.FindServerData(this.m_configInfo.MqttLoadBalanceInfo.MqttLoadBalanceServerName)
	if err != nil {
		log.Println("find server data from service discovery error, err: ", err)
		return nil, nil, err
	}
	var host *string = nil
	var ip string
	var port int
	var userName string
	var userPwd string
	var serverUniqueCode string
	if nodeItem.NormalNodes == nil || len(*nodeItem.NormalNodes) == 0 {
		if nodeItem.MasterNode == nil {
			log.Println("normal node or master node both is not exist")
			return nil, nil, errors.New("normal and master node is not exist")
		} else {
			ip = nodeItem.MasterNode.ServerIp
			port = nodeItem.MasterNode.ServerPort
			userName = nodeItem.MasterNode.UserName
			userPwd = nodeItem.MasterNode.UserPwd
			serverUniqueCode = nodeItem.MasterNode.ServerUniqueCode
			host = this.joinHost(&nodeItem.MasterNode.ServerIp, nodeItem.MasterNode.ServerPort)
		}
	} else {
		if this.m_algorithm == nil {
			log.Println("normal node algorithm is nil")
			return nil, nil, errors.New("normal node algorithm is nil")
		}
		data, err := this.m_algorithm.Get(this.m_configInfo.MqttLoadBalanceInfo.MqttLoadBalanceServerName, nil)
		if err != nil {
			log.Println("get normal node error, err: ", err)
			return nil, nil, err
		}
		ip = data.ServerIp
		port = data.ServerPort
		userName = data.UserName
		userPwd = data.UserPwd
		serverUniqueCode = data.ServerUniqueCode
		host = this.joinHost(&data.ServerIp, data.ServerPort)
	}
	var mqttComm mqtt_comm.CMqttComm
	value, ok := this.m_mqttCommMap.Load(*host)
	if !ok {
		// not exist
		mqttComm = this.connectBroker(&ip, port, &userName, &userPwd)
		this.m_mqttCommMap.Store(*host, mqttComm)
	} else {
		// exist
		mqttComm = value.(mqtt_comm.CMqttComm)
		if mqttComm.IsConnect() == false {
			mqttComm = nil
			mqttComm = this.connectBroker(&ip, port, &userName, &userPwd)
			this.m_mqttCommMap.Store(*host, mqttComm)
		}
	}
	return mqttComm, &url.CUrlMaker{
		ServerUniqueCode: &serverUniqueCode,
	}, nil
}

func (this *CClient) connectBroker(ip *string, port int, userName *string, userPwd *string) mqtt_comm.CMqttComm {
	mqttComm := mqtt_comm.NewMqttComm(this.m_configInfo.MqttLoadBalanceInfo.MqttServerName, this.m_configInfo.MqttLoadBalanceInfo.MqttServerVersion, this.m_configInfo.MqttLoadBalanceInfo.MqttServerRecvQos)
	mqttComm.SetMessageBus(*ip, port, *userName, *userPwd)
	mqttComm.Connect(false)
	for _, item := range this.m_subscribeInfo {
		mqttComm.Subscribe(item.action, item.topic, item.qos, item.handler, item.userData)
	}
	return mqttComm
}

func (this *CClient) JoinTopic(topic string, serverUniqueCode *string) *string {
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

func (this *CClient) Subscribe(action string, topic string, qos int, handler mqtt_comm.CHandler, user interface{}) error {
	info := CSubscribeInfo{
		action:   action,
		topic:    topic,
		qos:      qos,
		handler:  handler,
		userData: user,
	}
	this.m_subscribeInfo = append(this.m_subscribeInfo, &info)
	return nil
}

func (this *CClient) joinHost(host *string, port int) *string {
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
		mqttLoadBalanceInfo.ServerMode,
		&conns,
		mqttLoadBalanceInfo.PathPrefix,
		10,
	)
	select {
	case <-connChan:
		break
	}
	this.m_algorithm = this.m_balance.GetNormalNodeAlgorithm(mqttLoadBalanceInfo.NormalNodeAlgorithm)
	if this.m_algorithm == nil {
		log.Fatalln("get normalenode algorithm error")
		return
	}
}
