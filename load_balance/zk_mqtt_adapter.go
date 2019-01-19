package load_balance

import (
	proto "../common_proto"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MwlLj/mqtt_comm"
	MQTT "github.com/eclipse/paho.mqtt.golang"
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
	m_mqTopicBrokerInfoList   []CMqTopicBrokerInfo
	m_mqTopicBrokerInfoMap    map[string]CMqTopicBrokerInfo
	m_mqTopicBrokerInfosMutex sync.Mutex
	m_mqConnectMap            map[*CMqTopicBrokerInfo]mqtt_comm.CMqttComm
	m_isRouterByTopic         bool
}

func (this *CZkMqttAdapter) init(conns *[]proto.CConnectProperty, pathPrefix string, connTimeoutS int) (<-chan bool, error) {
	this.m_isRouterByTopic = true
	return this.CZkAdapter.init(conns, pathPrefix, connTimeoutS)
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

func (this *CZkMqttAdapter) Run() error {
	if this.m_isRouterByTopic == true {
		for _, v := range this.m_mqTopicBrokerInfoMap {
			mqttComm := mqtt_comm.NewMqttComm("zk_mqtt_loadblance", "1.0", 0)
			if mqttComm == nil {
				return errors.New("new mqttcomm error")
			}
			mqttComm.SetMessageBus(v.info.Host, v.info.Port, v.info.UserName, v.info.UserPwd)
			// err := mqttComm.Subscribe(mqtt_comm.POST, "#", 0, &handlers.CPostUserHandle{}, this)
			// if err != nil {
			// 	return err
			// }
			mqttComm.Connect(false)
			this.m_mqConnectMap[&v] = mqttComm
		}
	} else {
	}
	return nil
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

func (this *CZkMqttAdapter) onSubscribeMessage(client MQTT.Client, message MQTT.Message) {
	topic := message.Topic()
	brokerInfo, err := this.findBroker(&topic)
	if err != nil {
		return
	}
	_, ok := this.m_mqConnectMap[brokerInfo]
	if !ok {
		return
	}
}

func (this *CZkMqttAdapter) GetNormalNodeAlgorithm(algorithm string) INormalNodeAlgorithm {
	if algorithm == AlgorithmRoundRobin {
		alg := CHttpRoundRobin{}
		err := alg.init(this)
		if err != nil {
			return nil
		}
		return &alg
	} else if algorithm == AlgorithmWeightRoundRobin {
		alg := CHttpWeightRoundRobin{}
		err := alg.init(this)
		if err != nil {
			return nil
		}
		return &alg
	} else if algorithm == AlgorithmRandom {
		alg := CHttpRandom{}
		err := alg.init(this)
		if err != nil {
			return nil
		}
		return &alg
	} else if algorithm == AlgorithmWeightRandom {
		alg := CHttpWeightRandom{}
		err := alg.init(this)
		if err != nil {
			return nil
		}
		return &alg
	} else if algorithm == AlgorithmIpHash {
		alg := CHttpIpHash{}
		err := alg.init(this)
		if err != nil {
			return nil
		}
		return &alg
	} else if algorithm == AlgorithmUrlHash {
		alg := CHttpUrlHash{}
		err := alg.init(this)
		if err != nil {
			return nil
		}
		return &alg
	} else if algorithm == AlgorithmLeastConnect {
		alg := CHttpLeastConnections{}
		err := alg.init(this)
		if err != nil {
			return nil
		}
		return &alg
	} else {
		return nil
	}
}
