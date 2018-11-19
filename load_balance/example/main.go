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
	bls, connChan := bl.New(bl.ServerModeZookeeper, &conns, "micro-service", 10)
	select {
	case <-connChan:
		break
	}
	data, err := bls.GetMasterNode("testserver")
	if err != nil {
		fmt.Println("[ERROR] get master node error, ", err)
	} else {
		fmt.Println(data.ServerIp, data.ServerPort, data.ServerUniqueCode)
	}
	var _ = bls
	for {
		time.Sleep(100 * time.Millisecond)
	}
}
