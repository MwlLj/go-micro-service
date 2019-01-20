package load_balance

import (
	"fmt"
	"testing"
)

func TestConfigReaderRead(t *testing.T) {
	reader := CConfigReader{}
	path := "config.json"
	config, err := reader.Read(&path)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(config)
}

func TestFindRuleInfoByTopic(t *testing.T) {
	reader := CConfigReader{}
	path := "config.json"
	err := reader.Init(&path)
	if err != nil {
		fmt.Println(err)
		return
	}
	topic := "/login"
	info, isExist, err := reader.FindRuleInfoByTopic(&topic)
	if err != nil {
		fmt.Println(err)
		return
	}
	if isExist {
		fmt.Println(info.IsMaster, info.ObjServerName)
	}
}
