package load_balance

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	proto "github.com/MwlLj/go-micro-service/common_proto"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"net/http"
	"net/http/httputil"
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
	m_normalNodeAlgorithm INormalNodeAlgorithm
	m_configReader        *CConfigReader
}

func (this *CZkHttpAdapter) init(conns *[]proto.CConnectProperty, pathPrefix string, connTimeoutS int) (<-chan bool, error) {
	return this.CZkAdapter.init(conns, pathPrefix, connTimeoutS)
}

func (this *CZkHttpAdapter) onMessage(w http.ResponseWriter, r *http.Request) error {
	top := r.URL.String()
	ruleInfo, isFind, err := this.m_configReader.FindRuleInfoByTopic(&top)
	if err != nil {
		log.Println("find router rule by topic error, err: ", err)
		return err
	}
	var hostBuffer bytes.Buffer
	if isFind == false {
		log.Println("not find router rule by topic")
		return errors.New("rule is not match")
	} else {
		var nodeData *proto.CNodeData = nil
		if ruleInfo.IsMaster {
			nodeData, err = this.GetMasterNode(ruleInfo.ObjServerName)
			if err != nil {
				fmt.Println(err)
				return err
			}
		} else {
			nodeData, err = this.m_normalNodeAlgorithm.Get(ruleInfo.ObjServerName, nil)
			if err != nil {
				nodeData, err = this.GetMasterNode(ruleInfo.ObjServerName)
				if err != nil {
					fmt.Println(err)
					return err
				}
			}
		}
		hostBuffer.WriteString(nodeData.ServerIp)
		hostBuffer.WriteString(":")
		hostBuffer.WriteString(strconv.FormatInt(int64(nodeData.ServerPort), 10))
	}
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = hostBuffer.String()
		req.URL.Path = top
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(w, r)
	return nil
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

func (this *CZkHttpAdapter) SetNormalNodeAlgorithm(algorithm string) error {
	this.m_normalNodeAlgorithm = this.GetNormalNodeAlgorithm(algorithm)
	if this.m_normalNodeAlgorithm == nil {
		return errors.New("not support")
	}
	return nil
}

func (this *CZkHttpAdapter) Run(data interface{}) error {
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
	serverInfo := data.(CNetInfo)
	addrBuffer := bytes.Buffer{}
	addrBuffer.WriteString(serverInfo.Host)
	addrBuffer.WriteString(":")
	addrBuffer.WriteString(strconv.FormatInt(int64(serverInfo.Port), 10))
	http.ListenAndServe(addrBuffer.String(), &CHttpMux{
		m_handler: func(w http.ResponseWriter, r *http.Request) {
			this.onMessage(w, r)
		},
	})
	return err
}

type CHttpMux struct {
	m_handler func(w http.ResponseWriter, r *http.Request)
}

func (this *CHttpMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	this.m_handler(w, r)
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
