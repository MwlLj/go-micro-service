package load_balance

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/MwlLj/gotools/configs"
	"regexp"
)

type CConfigReader struct {
	configs.CConfigBase
	m_configInfo *CConfigInfo
}

func (this *CConfigReader) Init(info *CConfigInfo) error {
	var err error = nil
	/*
		this.m_configInfo, err = this.Read(path)
		if err != nil {
			return err
		}
	*/
	this.m_configInfo = info
	return err
}

func (this *CConfigReader) FindRuleInfoByTopic(topic *string) (*CRuleInfo, bool, error) {
	if this.m_configInfo == nil {
		return nil, false, errors.New("not init")
	}
	for k, v := range this.m_configInfo.Rules {
		ok, err := regexp.Match(k, []byte(*topic))
		if err != nil {
			return nil, false, err
		}
		if ok {
			return v, true, nil
		}
	}
	return nil, false, nil
}

func (this *CConfigReader) Read(path *string) (*CConfigInfo, error) {
	rules := map[string]*CRuleInfo{}
	topic := ".*"
	rule := CRuleInfo{
		ObjServerName: "",
		IsMaster:      false,
	}
	rules[topic] = &rule
	topic = "/login"
	rule = CRuleInfo{
		ObjServerName: "",
		IsMaster:      false,
	}
	rules[topic] = &rule
	configInfo := CConfigInfo{
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
	outputConfigInfo := CConfigInfo{}
	err = json.Unmarshal([]byte(value), &outputConfigInfo)
	if err != nil {
		return nil, err
	}
	return &outputConfigInfo, nil
}
