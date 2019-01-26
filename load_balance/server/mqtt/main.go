package main

import (
	bl "../.."
	proto "../../../common_proto"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/MwlLj/gotools/configs"
)

var _ = fmt.Println

type CConfigReader struct {
	configs.CConfigBase
}

func (this *CConfigReader) ReadConfig(path *string) (*bl.CConfigInfo, error) {
	rules := map[string]*bl.CRuleInfo{}
	topic := ".*"
	rule := bl.CRuleInfo{
		ObjServerName: "",
		IsMaster:      false,
	}
	rules[topic] = &rule
	topic = "/login"
	rule = bl.CRuleInfo{
		ObjServerName: "",
		IsMaster:      false,
	}
	rules[topic] = &rule
	configInfo := bl.CConfigInfo{
		Rules: rules,
	}
	// parse default
	b, err := json.Marshal(&configInfo)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "\t")
	if err != nil {
		return nil, err
	}
	value, err := this.Load(*path, out.String())
	if err != nil {
		return nil, err
	}
	// parse output
	outputConfigInfo := bl.CConfigInfo{}
	err = json.Unmarshal([]byte(value), &outputConfigInfo)
	if err != nil {
		return nil, err
	}
	return &outputConfigInfo, nil
}

func mqttTest() {
	var conns []proto.CConnectProperty
	conns = append(conns, proto.CConnectProperty{ServerHost: "127.0.0.1", ServerPort: 2182, ServiceId: "server_1"})
	bls, connChan := bl.New(bl.ServerModeZookeeperMqtt, &conns, "micro-service", 10)
	// algorithm := bls.GetNormalNodeAlgorithm(bl.AlgorithmWeightRoundRobin)
	// algorithm := bls.GetNormalNodeAlgorithm(bl.AlgorithmRandom)
	// algorithm := bls.GetNormalNodeAlgorithm(bl.AlgorithmWeightRandom)
	// algorithm := bls.GetNormalNodeAlgorithm(bl.AlgorithmIpHash)
	// algorithm := bls.GetNormalNodeAlgorithm(bl.AlgorithmUrlHash)
	// algorithm := bls.GetNormalNodeAlgorithm(bl.AlgorithmLeastConnect)
	select {
	case <-connChan:
		break
	}
	err := bls.SetNormalNodeAlgorithm(bl.AlgorithmRoundRobin)
	if err != nil {
		fmt.Println(err)
		return
	}
	configReader := CConfigReader{}
	cfgPath := "mqtt_load_blance.cfg"
	cfgInfo, err := configReader.ReadConfig(&cfgPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	bls.SetConfigInfo(cfgInfo)
	bls.AddRecvNetInfo(nil, &bl.CNetInfo{
		Host: "localhost",
		Port: 51883,
	})
	err = bls.Run(bl.CNetInfo{
		Host: "localhost",
		Port: 51885,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	mqttTest()
}
