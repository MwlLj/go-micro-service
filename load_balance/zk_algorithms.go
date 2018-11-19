package load_balance

import (
	proto "../common_proto"
	"errors"
	"fmt"
)

var _ = fmt.Println
var _ = errors.New

type CRoundRobin struct {
	m_loadBlance ILoadBlance
	m_nodeLength int
	m_nodeIndex  int
}

func (this *CRoundRobin) Get(serverName string) (*proto.CNodeData, error) {
	item, err := this.m_loadBlance.findServerData(serverName)
	if err != nil {
		fmt.Println("[ERROR] server not find")
		return nil, err
	}
	length := len(*item.normalNodes)
	if length != this.m_nodeLength || this.m_nodeLength == 0 {
		this.m_nodeLength = length
		this.m_nodeIndex = 0
	}
	data := (*item.normalNodes)[this.m_nodeIndex]
	this.m_nodeIndex += 1
	if this.m_nodeIndex > this.m_nodeLength-1 {
		this.m_nodeIndex = 0
	}
	return &data, nil
}

type CWeightRoundRobin struct {
	m_loadBlance ILoadBlance
}

func (this *CWeightRoundRobin) Get(serverName string) (*proto.CNodeData, error) {
	return nil, nil
}

type CRandom struct {
	m_loadBlance ILoadBlance
}

func (this *CRandom) Get(serverName string) (*proto.CNodeData, error) {
	return nil, nil
}

type CWeightRandom struct {
	m_loadBlance ILoadBlance
}

func (this *CWeightRandom) Get(serverName string) (*proto.CNodeData, error) {
	return nil, nil
}

type CIpHash struct {
	m_loadBlance ILoadBlance
}

func (this *CIpHash) Get(serverName string) (*proto.CNodeData, error) {
	return nil, nil
}

type CUrlHash struct {
	m_loadBlance ILoadBlance
}

func (this *CUrlHash) Get(serverName string) (*proto.CNodeData, error) {
	return nil, nil
}

type CLeastConnections struct {
	m_loadBlance ILoadBlance
}

func (this *CLeastConnections) Get(serverName string) (*proto.CNodeData, error) {
	return nil, nil
}
