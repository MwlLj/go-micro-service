package main

import (
	s ".."
	"fmt"
)

func main() {
	var conns []s.CConnectProperty
	conns = append(conns, s.CConnectProperty{ServerHost: "127.0.0.1", ServerPort: 2181, ServiceId: "server_1"})
	sds := s.New(&s.CInitProperty{ServerMode: s.ServerModeZookeeper, Conns: conns, ConnTimeoutS: 10})
	err := sds.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("success")
}
