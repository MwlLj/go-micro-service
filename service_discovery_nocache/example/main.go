package main

import (
	"fmt"
	proto "github.com/MwlLj/go-micro-service/common_proto"
	s "github.com/MwlLj/go-micro-service/service_discovery_nocache"
	"time"
)

var _ = fmt.Println

func main() {
	var conns []proto.CConnectProperty
	conns = append(conns, proto.CConnectProperty{ServerHost: "127.0.0.1", ServerPort: 2182, ServiceId: "server_1"})
	sds := s.New(&s.CInitProperty{
		PathPrefix:   "micro-service",
		ServerMode:   s.ServerModeZookeeper,
		ServerName:   "testserver",
		NodeData:     proto.CNodeData{ServerIp: "127.0.0.1", ServerPort: 50000, ServerUniqueCode: "cacd3aa4-4eb8-4bf6-b967-fbcee5377992", Weight: 1},
		Conns:        conns,
		ConnTimeoutS: 10})
	var _ = sds
	for {
		time.Sleep(100 * time.Millisecond)
	}
}
