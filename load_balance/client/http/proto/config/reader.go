package config

import (
	"bytes"
	"encoding/json"
	"github.com/MwlLj/gotools/configs"
	"log"
)

type CReader struct {
	configs.CConfigBase
}

func (this *CReader) Read(path string) (*CRemoteProtoConfigInfo, error) {
	// default value
	configInfo := CRemoteProtoConfigInfo{}
	// start rule
	startRule := CRemoteProtoStartProto{}
	startRule.Url = "/server/init"
	configInfo.StartRule = startRule
	// get connect rule
	getConnectRule := CRemoteProtoGetConnectProto{}
	getConnectRule.Url = "/connect/loadbalance"
	configInfo.GetConnectRule = getConnectRule
	// register service rule
	registerServiceRule := CRemoteProtoRegisterServiceProto{}
	registerServiceRule.Url = "/service/register"
	configInfo.RegisterServiceRule = registerServiceRule
	// net info
	netInfo := CRemoteListenNetInfo{}
	netInfo.Host = "localhost"
	netInfo.Port = 15000
	configInfo.Net = netInfo
	// encoder
	b, err := json.Marshal(&configInfo)
	if err != nil {
		log.Fatalln("read config encoder error, err: ", err)
		return nil, err
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "\t")
	if err != nil {
		log.Fatalln("config file format change error, err: ", err)
		return nil, err
	}
	value, err := this.Load(path, out.String())
	if err != nil {
		log.Fatalln("read config error, err: ", err)
		return nil, err
	}
	// parse output
	outputConfig := CRemoteProtoConfigInfo{}
	err = json.Unmarshal([]byte(value), &outputConfig)
	if err != nil {
		log.Fatalln("parer config error, err: ", err)
		return nil, err
	}
	return &outputConfig, nil
}
