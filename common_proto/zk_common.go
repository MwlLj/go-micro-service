package common_proto

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

const MasterNode string = "master"
const NormalNode string = "normal"

type CConnInfo struct {
	Host string
}

type CNodeData struct {
	ServerIp         string `json:"serverip"`
	ServerPort       int    `json:"serverport"`
	ServerUniqueCode string `json:"serveruniquecode"`
	Weight           int    `json:"weight"`
}

type CConnectProperty struct {
	ServerHost string
	ServerPort int
	ServiceId  string
}

type CZkCommon struct {
	m_connMap sync.Map
}

func (*CZkCommon) JoinPathPrefix(prefix *string, serverName *string) *string {
	var path string
	if prefix != nil {
		if serverName != nil {
			path = strings.Join([]string{*prefix, *serverName}, "/")
		} else {
			path = *prefix
		}
	} else {
		if serverName != nil {
			path = *serverName
		} else {
			path = ""
		}
	}
	path = strings.Join([]string{"/", path}, "")
	return &path
}

func (*CZkCommon) GetParentNode(path string) (isRoot bool, parent string) {
	li := strings.Split(path, "/")
	length := len(li)
	if length == 1 {
		return true, li[0]
	}
	return false, strings.Join(li[:length-1], "/")
}

func (*CZkCommon) JoinHost(ip string, port int) *string {
	host := strings.Join([]string{ip, strconv.FormatInt(int64(port), 10)}, ":")
	return &host
}

func (this *CZkCommon) ToHosts() *[]string {
	var hosts []string
	f := func(k, v interface{}) bool {
		info := v.(CConnInfo)
		hosts = append(hosts, info.Host)
		return true
	}
	this.m_connMap.Range(f)
	return &hosts
}

func (this *CZkCommon) AddConnProperty(conn *CConnectProperty) error {
	info := CConnInfo{}
	info.Host = *this.JoinHost(conn.ServerHost, conn.ServerPort)
	this.m_connMap.Store(conn.ServiceId, info)
	return nil
}

func (this *CZkCommon) UpdateConnProperty(conn *CConnectProperty) error {
	info := CConnInfo{}
	info.Host = *this.JoinHost(conn.ServerHost, conn.ServerPort)
	this.m_connMap.Store(conn.ServiceId, info)
	return nil
}

func (this *CZkCommon) DeleteConnProperty(serviceId *string) error {
	this.m_connMap.Delete(serviceId)
	return nil
}

func (this *CZkCommon) SplitePath(path string) (prefix, serverName, nodeName *string, e error) {
	paths := strings.Split(path, "/")
	length := len(paths)
	if length < 4 {
		return nil, nil, nil, errors.New("[ERROR] path rule is error")
	}
	*prefix = strings.Join(paths[1:length-2], "/")
	*serverName = paths[length-2]
	*nodeName = paths[length-1]
	return prefix, serverName, nodeName, nil
}
