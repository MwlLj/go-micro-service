package service_discovery_single

var (
	ServerModeZookeeper = "zookeeper"
)

type CInitProperty struct {
	// zookeeper ...
	ServerMode string
}

type CConnectProperty struct {
	ZkServerHost string
	ZkServerPort int
	ConnTimeoutS int
	ServiceName  string
}
