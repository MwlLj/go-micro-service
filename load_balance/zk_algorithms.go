package load_balance

import (
	proto "../common_proto"
)

type CRoundRobin struct {
	m_loadBlance ILoadBlance
}

func (this *CRoundRobin) Get(serverName string) (*proto.CNodeData, error) {
	return nil, nil
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
