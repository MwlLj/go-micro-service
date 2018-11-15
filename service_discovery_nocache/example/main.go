package main

import (
	s ".."
	"fmt"
	"time"
)

func main() {
	var conns []s.CConnectProperty
	conns = append(conns, s.CConnectProperty{ServerHost: "127.0.0.1", ServerPort: 2182, ServiceId: "server_1"})
	sds := s.New(&s.CInitProperty{
		PathPrefix:   "mico-service",
		ServerMode:   s.ServerModeZookeeper,
		ServerName:   "testserver",
		Conns:        conns,
		ConnTimeoutS: 10})
	err := sds.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("success")
	for {
		time.Sleep(100 * time.Millisecond)
	}
}
