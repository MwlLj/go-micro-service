package load_balance

import (
	proto "../common_proto"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MwlLj/mqtt_comm"
	"github.com/samuel/go-zookeeper/zk"
	"strconv"
	"strings"
	"sync"
	"time"
)

var _ = fmt.Println
var _ = errors.New("")
var _ = strings.Join
var _ = json.Unmarshal
var _ = zk.Connect
var _ = strconv.FormatInt
var _ = time.Sleep
var _ = sync.NewCond

type CMqTopicBrokerInfo struct {
	info  *CNetInfo
	topic *string
}

type CZkMqttAdapter struct {
	CZkAdapter
	m_normalNodeAlgorithm     INormalNodeAlgorithm
	m_transmitTimeoutS        int
	m_configFilePath          *string
	m_configReader            *CConfigReader
	m_mqttComm                mqtt_comm.CMqttComm
	m_mqTopicBrokerInfoList   []CMqTopicBrokerInfo
	m_mqTopicBrokerInfoMap    map[string]CMqTopicBrokerInfo
	m_mqTopicBrokerInfosMutex sync.Mutex
	m_mqConnectMap            map[*CMqTopicBrokerInfo]mqtt_comm.CMqttComm
	m_mqConnectMapMutex       sync.Mutex
	m_isRouterByTopic         bool
}

func (this *CZkMqttAdapter) init(conns *[]proto.CConnectProperty, pathPrefix string, connTimeoutS int) (<-chan bool, error) {
	this.m_isRouterByTopic = true
	this.m_transmitTimeoutS = 60
	path := "rule-config.json"
	this.m_configFilePath = &path
	return this.CZkAdapter.init(conns, pathPrefix, connTimeoutS)
}

func (this *CZkMqttAdapter) SetTransmitTimeoutS(s int) {
	this.m_transmitTimeoutS = s
}

func (this *CZkMqttAdapter) onMessage(topic *string, action *string, request *string, qos int) (*string, error) {
	brokerInfo, err := this.findBroker(topic)
	if err != nil {
		return nil, err
	}
	mqttComm, ok := this.m_mqConnectMap[brokerInfo]
	if !ok {
		return nil, errors.New("not found")
	}
	ruleInfo, isFind, err := this.m_configReader.FindRuleInfoByTopic(topic)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	if isFind == false {
		return nil, errors.New("rule is not match")
	} else {
		var nodeData *proto.CNodeData = nil
		if ruleInfo.IsMaster {
			nodeData, err = this.GetMasterNode(ruleInfo.ObjServerName)
			if err != nil {
				return nil, err
			}
		} else {
			nodeData, err = this.m_normalNodeAlgorithm.Get(ruleInfo.ObjServerName, nil)
			if err != nil {
				return nil, err
			}
		}
		serverUuid := nodeData.ServerUniqueCode
		buffer.WriteString(*topic)
		buffer.WriteString("/")
		buffer.WriteString(serverUuid)
	}
	response, err := mqttComm.Send(*action, buffer.String(), *request, qos, this.m_transmitTimeoutS)
	return &response, err
}

func (this *CZkMqttAdapter) AddRecvNetInfo(topic *string, info *CNetInfo) {
	if topic == nil || *topic == "" {
		this.m_isRouterByTopic = false
	}
	brokerInfo := CMqTopicBrokerInfo{
		info:  info,
		topic: topic,
	}
	this.m_mqTopicBrokerInfosMutex.Lock()
	this.m_mqTopicBrokerInfoList = append(this.m_mqTopicBrokerInfoList, brokerInfo)
	this.m_mqTopicBrokerInfosMutex.Unlock()
	if topic != nil && *topic != "" {
		this.m_mqTopicBrokerInfosMutex.Lock()
		this.m_mqTopicBrokerInfoMap[*topic] = brokerInfo
		this.m_mqTopicBrokerInfosMutex.Unlock()
	}
}

func (this *CZkMqttAdapter) Run(data interface{}) error {
	if this.m_configReader == nil {
		return errors.New("please call SetConfigFilePath")
	}
	if this.m_normalNodeAlgorithm == nil {
		return errors.New("please call SetNormalNodeAlgorithm")
	}
	var err error = nil
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("data struct error, please use CNetInfo")
		}
	}()
	// connect recvs brokers
	connInner := func(info *CMqTopicBrokerInfo) error {
		mqttComm := mqtt_comm.NewMqttComm("zk_mqtt_loadblance", "1.0", 0)
		if mqttComm == nil {
			return errors.New("new mqttcomm error")
		}
		mqttComm.SetMessageBus(info.info.Host, info.info.Port, info.info.UserName, info.info.UserPwd)
		mqttComm.Connect(false)
		this.m_mqConnectMapMutex.Lock()
		this.m_mqConnectMap[info] = mqttComm
		this.m_mqConnectMapMutex.Unlock()
		return nil
	}
	if this.m_isRouterByTopic == true {
		for _, v := range this.m_mqTopicBrokerInfoMap {
			err = connInner(&v)
			if err != nil {
				return err
			}
		}
	} else {
		for _, v := range this.m_mqTopicBrokerInfoList {
			err = connInner(&v)
			if err != nil {
				return err
			}
		}
	}
	// connect broker
	brokerNetInfo := data.(CNetInfo)
	this.m_mqttComm = mqtt_comm.NewMqttComm("cfgs", "1.0", 0)
	this.m_mqttComm.SetMessageBus(brokerNetInfo.Host, brokerNetInfo.Port, brokerNetInfo.UserName, brokerNetInfo.UserPwd)
	this.m_mqttComm.SubscribeAll(0, &CRequestHandler{}, this)
	this.m_mqttComm.Connect(true)
	return err
}

func (this *CZkMqttAdapter) findBroker(topic *string) (*CMqTopicBrokerInfo, error) {
	if this.m_isRouterByTopic == true {
		v, ok := this.m_mqTopicBrokerInfoMap[*topic]
		if ok {
			return &v, nil
		}
	} else {
		randValue := randomInt(0, len(this.m_mqTopicBrokerInfoList))
		return &this.m_mqTopicBrokerInfoList[randValue], nil
	}
	return nil, errors.New("not found")
}

func (this *CZkMqttAdapter) SetNormalNodeAlgorithm(algorithm string) error {
	this.m_normalNodeAlgorithm = this.GetNormalNodeAlgorithm(algorithm)
	if this.m_normalNodeAlgorithm == nil {
		return errors.New("not support")
	}
	return nil
}

func (this *CZkMqttAdapter) SetConfigFilePath(path *string) error {
	if path != nil {
		this.m_configFilePath = path
	}
	var err error = nil
	if this.m_configReader == nil {
		this.m_configReader = &CConfigReader{}
		err = this.m_configReader.Init(this.m_configFilePath)
		if err != nil {
			this.m_configReader = nil
		}
	}
	return err
}

func (this *CZkMqttAdapter) GetNormalNodeAlgorithm(algorithm string) INormalNodeAlgorithm {
	if algorithm == AlgorithmRoundRobin {
		alg := CMqttRoundRobin{}
		err := alg.init(this)
		if err != nil {
			return nil
		}
		return &alg
	} else if algorithm == AlgorithmWeightRoundRobin {
		alg := CMqttWeightRoundRobin{}
		err := alg.init(this)
		if err != nil {
			return nil
		}
		return &alg
	} else if algorithm == AlgorithmRandom {
		alg := CMqttRandom{}
		err := alg.init(this)
		if err != nil {
			return nil
		}
		return &alg
	} else if algorithm == AlgorithmWeightRandom {
		alg := CMqttWeightRandom{}
		err := alg.init(this)
		if err != nil {
			return nil
		}
		return &alg
	} else if algorithm == AlgorithmIpHash {
		alg := CMqttIpHash{}
		err := alg.init(this)
		if err != nil {
			return nil
		}
		return &alg
	} else if algorithm == AlgorithmUrlHash {
		alg := CMqttUrlHash{}
		err := alg.init(this)
		if err != nil {
			return nil
		}
		return &alg
	} else if algorithm == AlgorithmLeastConnect {
		alg := CMqttLeastConnections{}
		err := alg.init(this)
		if err != nil {
			return nil
		}
		return &alg
	} else {
		return nil
	}
}

type CRequestHandler struct {
}

func (this *CRequestHandler) Handle(topic *string, action *string, request *string, qos int, mc mqtt_comm.CMqttComm, user interface{}) (*string, error) {
	adapter := user.(CZkMqttAdapter)
	return adapter.onMessage(topic, action, request, qos)
}
