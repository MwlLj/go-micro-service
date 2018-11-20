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
	// algorithm := bls.GetNormalNodeAlgorithm(bl.AlgorithmRoundRobin)
	// algorithm := bls.GetNormalNodeAlgorithm(bl.AlgorithmWeightRoundRobin)
	// algorithm := bls.GetNormalNodeAlgorithm(bl.AlgorithmRandom)
	algorithm := bls.GetNormalNodeAlgorithm(bl.AlgorithmWeightRandom)
	select {
	case <-connChan:
		break
	}
	data, err := bls.GetMasterNode("testserver")
	if err != nil {
		fmt.Println("[ERROR] get master node error, ", err)
	} else {
		fmt.Println("------------master------------")
		fmt.Println(data.ServerIp, data.ServerPort, data.ServerUniqueCode, data.Weight)
		fmt.Println("------------master------------")
	}
	fmt.Println("------------normal------------")
	for i := 0; i < 20; i++ {
		data, err = algorithm.Get("testserver")
		if err != nil {
			fmt.Println("[ERROR] get normal node error, ", err)
		} else {
			fmt.Println(data.ServerIp, data.ServerPort, data.ServerUniqueCode, data.Weight)
		}
	}
	fmt.Println("------------normal------------")
	var _ = bls
	for {
		time.Sleep(100 * time.Millisecond)
	}
}
