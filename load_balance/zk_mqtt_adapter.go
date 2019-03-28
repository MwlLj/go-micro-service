package load_balance

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	proto "github.com/MwlLj/go-micro-service/common_proto"
	"github.com/MwlLj/mqtt_comm"
	"github.com/samuel/go-zookeeper/zk"
	"log"
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
	m_configReader            *CConfigReader
	m_mqttComm                mqtt_comm.CMqttComm
	m_mqTopicBrokerInfoList   []CMqTopicBrokerInfo
	m_mqTopicBrokerInfoMap    map[string]CMqTopicBrokerInfo
	m_mqTopicBrokerInfosMutex sync.Mutex
	m_mqConnectMap            sync.Map
	m_isRouterByTopic         bool
}

func (this *CZkMqttAdapter) init(conns *[]proto.CConnectProperty, pathPrefix string, connTimeoutS int) (<-chan bool, error) {
	this.m_isRouterByTopic = true
	this.m_transmitTimeoutS = 60
	this.m_mqTopicBrokerInfoMap = make(map[string]CMqTopicBrokerInfo)
	// this.m_mqConnectMap = make(map[string]mqtt_comm.CMqttComm)
	return this.CZkAdapter.init(conns, pathPrefix, connTimeoutS)
}

func (this *CZkMqttAdapter) brokerHostJoin(host *string, port int) *string {
	var buffer bytes.Buffer
	buffer.WriteString(*host)
	buffer.WriteString(strconv.FormatInt(int64(port), 10))
	s := buffer.String()
	return &s
}

func (this *CZkMqttAdapter) TopicJoin(topic *string, serverUniqueCode *string) *string {
	var buffer bytes.Buffer
	buffer.WriteString(*topic)
	bTopic := []byte(*topic)
	if bTopic[len(bTopic)-1] != '/' {
		buffer.WriteString("/")
	}
	buffer.WriteString(*serverUniqueCode)
	top := buffer.String()
	return &top
}

func (this *CZkMqttAdapter) SetTransmitTimeoutS(s int) {
	this.m_transmitTimeoutS = s
}

func (this *CZkMqttAdapter) delTopicPrefix(topic *string) *string {
	index := strings.Index(*topic, "/")
	bTopic := []byte(*topic)
	s := string(bTopic[(index + 1):len(bTopic)])
	return &s
}

func (this *CZkMqttAdapter) connectBroker(brokerInfo *CMqTopicBrokerInfo) mqtt_comm.CMqttComm {
	mqttComm := mqtt_comm.NewMqttComm("zk-mqtt-loadblance", "1.0", 0)
	mqttComm.SetMessageBus(
		brokerInfo.info.Host,
		brokerInfo.info.Port,
		brokerInfo.info.UserName,
		brokerInfo.info.UserPwd,
	)
	mqttComm.Connect(false)
	return mqttComm
}

func (this *CZkMqttAdapter) findMqttConnect(brokerInfo *CMqTopicBrokerInfo) (mqtt_comm.CMqttComm, error) {
	brokerHash := this.brokerHostJoin(&brokerInfo.info.Host, brokerInfo.info.Port)
	value, ok := this.m_mqConnectMap.Load(brokerHash)
	var mqttComm mqtt_comm.CMqttComm = nil
	if !ok {
		// not exist
		mqttComm = this.connectBroker(brokerInfo)
		this.m_mqConnectMap.Store(brokerHash, mqttComm)
	} else {
		// exist
		mqttComm = value.(mqtt_comm.CMqttComm)
		if mqttComm.IsConnect() == false {
			mqttComm = nil
			mqttComm = this.connectBroker(brokerInfo)
		}
	}
	return mqttComm, nil
}

func (this *CZkMqttAdapter) onMessage(topic *string, action *string, request *string, qos int) (*string, error) {
	brokerInfo, err := this.findBroker(topic)
	if err != nil {
		log.Println("find recver broker by topic error, err: ", err)
		return nil, err
	}
	mqttComm, err := this.findMqttConnect(brokerInfo)
	if err != nil {
		log.Println("connect recver broker error")
		return nil, err
	}
	top := this.delTopicPrefix(topic)
	ruleInfo, isFind, err := this.m_configReader.FindRuleInfoByTopic(top)
	if err != nil {
		log.Println("find router rule by topic error, err: ", err)
		return nil, err
	}
	var buffer bytes.Buffer
	if isFind == false {
		log.Println("not find router rule by topic")
		return nil, errors.New("rule is not match")
	} else {
		var nodeData *proto.CNodeData = nil
		if ruleInfo.IsMaster {
			nodeData, err = this.GetMasterNode(ruleInfo.ObjServerName)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
		} else {
			nodeData, err = this.m_normalNodeAlgorithm.Get(ruleInfo.ObjServerName, nil)
			if err != nil {
				nodeData, err = this.GetMasterNode(ruleInfo.ObjServerName)
				if err != nil {
					fmt.Println(err)
					return nil, err
				}
			}
		}
		serverUuid := nodeData.ServerUniqueCode
		buffer.WriteString(*top)
		// buffer.WriteString("/")
		buffer.WriteString(serverUuid)
	}
	response, err := mqttComm.Send(*action, buffer.String(), *request, true, qos, this.m_transmitTimeoutS)
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
	// connect broker
	brokerNetInfo := data.(CNetInfo)
	this.m_mqttComm = mqtt_comm.NewMqttComm("", "", 0)
	this.m_mqttComm.SetMessageBus(brokerNetInfo.Host, brokerNetInfo.Port, brokerNetInfo.UserName, brokerNetInfo.UserPwd)
	this.m_mqttComm.SubscribeAll(&brokerNetInfo.ExtraField, 0, &CRequestHandler{}, this)
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

func (this *CZkMqttAdapter) SetConfigInfo(info *CConfigInfo) error {
	var err error = nil
	if this.m_configReader == nil {
		this.m_configReader = &CConfigReader{}
		err = this.m_configReader.Init(info)
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
	adapter := user.(*CZkMqttAdapter)
	return adapter.onMessage(topic, action, request, qos)
}
