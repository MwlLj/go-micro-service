package main

import (
	s ".."
	proto "../../common_proto"
	"fmt"
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
