package main

import (
	s ".."
)

func main() {
	sds := s.New(&s.CInitProperty{ServerMode: s.ServerModeZookeeper})
	sds.Connect(&s.CConnectProperty{ZkServerHost: "192.168.9.15", ZkServerPort: 2181, ServiceName: "svrtest"})
}
