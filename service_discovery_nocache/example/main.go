package main

import (
	s ".."
	"fmt"
)

func main() {
	var nets []s.CNet
	nets = append(nets, s.CNet{ServerHost: "127.0.0.1", ServerPort: 2181})
	sds := s.New(&s.CInitProperty{ServerMode: s.ServerModeZookeeper, ConnProperty: s.CConnectProperty{nets}})
	err := sds.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("success")
}
