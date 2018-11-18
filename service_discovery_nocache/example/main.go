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
		PathPrefix:       "mico-service",
		ServerMode:       s.ServerModeZookeeper,
		ServerName:       "testserver",
		ServerUniqueCode: "0B5398177EBA429898F68AF13909920E",
		NodePayload:      `{"ip":"192.168.9.15","port":50000}`,
		Conns:            conns,
		ConnTimeoutS:     10})
	var _ = sds
	for {
		time.Sleep(100 * time.Millisecond)
	}
}
