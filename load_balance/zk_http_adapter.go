package load_balance

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	proto "github.com/MwlLj/go-micro-service/common_proto"
	"github.com/samuel/go-zookeeper/zk"
	"net/http"
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

type CZkHttpAdapter struct {
	CZkAdapter
	m_configReader *CConfigReader
}

func (this *CZkHttpAdapter) init(conns *[]proto.CConnectProperty, pathPrefix string, connTimeoutS int) (<-chan bool, error) {
	return this.CZkAdapter.init(conns, pathPrefix, connTimeoutS)
}

func (this *CZkMqttAdapter) onMessage(w http.ResponseWriter, r *http.Request) {
	top := r.URL.String()
	ruleInfo, isFind, err := this.m_configReader.FindRuleInfoByTopic(&top)
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
	response, err := mqttComm.Send(*action, buffer.String(), *request, qos, this.m_transmitTimeoutS)
}

func (this *CZkHttpAdapter) SetConfigInfo(info *CConfigInfo) error {
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

func (this *CZkHttpAdapter) Run(data interface{}) error {
	var err error = nil
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("data struct error, please use CNetInfo")
		}
	}()
	// connect broker
	serverInfo := data.(CNetInfo)
	addrBuffer := bytes.Buffer{}
	addrBuffer.WriteString(serverInfo.Host)
	addrBuffer.WriteString(":")
	addrBuffer.WriteString(strconv.FormatInt(int64(serverInfo.Port), 10))
	http.ListenAndServe(addrBuffer.String(), func(w http.ResponseWriter, r *http.Request) {
		this.onMessage(w, r)
	})
	return err
}

func (this *CZkHttpAdapter) GetNormalNodeAlgorithm(algorithm string) INormalNodeAlgorithm {
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
