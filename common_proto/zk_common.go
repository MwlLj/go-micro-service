package common_proto

import (
	"strconv"
	"strings"
	"sync"
)

type CConnInfo struct {
	Host string
}

func JoinPathPrefix(prefix *string, serverName *string) *string {
	var path string
	if *prefix != "" {
		path = strings.Join([]string{*prefix, *serverName}, "/")
	} else {
		path = *serverName
	}
	path = strings.Join([]string{"/", path}, "")
	return &path
}

func GetParentNode(path string) (isRoot bool, parent string) {
	li := strings.Split(path, "/")
	length := len(li)
	if length == 1 {
		return true, li[0]
	}
	return false, strings.Join(li[:length-1], "/")
}

func JoinHost(ip string, port int) *string {
	host := strings.Join([]string{ip, strconv.FormatInt(int64(port), 10)}, ":")
	return &host
}

func ToHosts(connMap *sync.Map) *[]string {
	var hosts []string
	f := func(k, v interface{}) bool {
		info := v.(CConnInfo)
		hosts = append(hosts, info.Host)
		return true
	}
	connMap.Range(f)
	return &hosts
}
