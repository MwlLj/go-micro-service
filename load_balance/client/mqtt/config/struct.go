package config

type CServiceDiscoveryNet struct {
	ServerHost string `json:"server-host"`
	ServerPort int    `json:"server-port"`
	ServiceId  string `json:"service-id"`
}

type CLoadBalanceInfo struct {
	PathPrefix                string                 `json:"path-prefix"`
	ServerMode                string                 `json:"server-mode"`
	NormalNodeAlgorithm       string                 `json:"normal-node-algorithm"`
	MqttLoadBalanceServerName string                 `json:"mqtt-loadbalance-servername"`
	Conns                     []CServiceDiscoveryNet `json:"conns"`
}

type CServiceInfo struct {
	PathPrefix            string                 `json:"path-prefix"`
	ServerMode            string                 `json:"server-mode"`
	ServerName            string                 `json:"server-name"`
	ServerVersion         string                 `json:"server-version"`
	ServerRecvQos         int                    `json:"server-recvqos"`
	ServerUniqueCode      string                 `json:"server-uniquecode"`
	Weight                int                    `json:"weight"`
	BrokerHost            string                 `json:"broker-host"`
	BrokerPort            int                    `json:"broker-port"`
	BrokerUserName        string                 `json:"broker-username"`
	BrokerUserPwd         string                 `json:"broker-userpwd"`
	ServiceDiscoveryConns []CServiceDiscoveryNet `json:"service-discovery-conns"`
}

type CConfigInfo struct {
	MqttLoadBalanceInfo CLoadBalanceInfo `json:"mqtt-loadbalance-info"`
	ServiceInfo         CServiceInfo     `json:"service-info"`
}
