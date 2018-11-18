package main

import (
	bl ".."
	proto "../../common_proto"
	"fmt"
	"time"
)

var _ = fmt.Println

func main() {
	var conns []proto.CConnectProperty
	conns = append(conns, proto.CConnectProperty{ServerHost: "127.0.0.1", ServerPort: 2182, ServiceId: "server_1"})
	bls := bl.New(bl.ServerModeNocacheZookeeper, &conns, "micro-service", 10)
	var _ = bls
	for {
		time.Sleep(100 * time.Millisecond)
	}
}
