package config

type CServiceDiscoryNet struct {
	ServerHost string `json:"server-host"`
	ServerPort int    `json:"server-port"`
	ServiceId  string `json:"service-id"`
}

type CLoadBalanceInfo struct {
	PathPrefix                string               `json:"path-prefix"`
	ServerMode                string               `json:"server-mode"`
	NormalNodeAlgorithm       string               `json:"normal-node-algorithm"`
	MqttLoadBalanceServerName string               `json:"mqtt-loadbalance-servername"`
	MqttServerName            string               `json:"mqtt-server-name"`
	MqttServerVersion         string               `json:"mqtt-server-version"`
	MqttServerRecvQos         int                  `json:"mqtt-server-recvqos"`
	Conns                     []CServiceDiscoryNet `json:"conns"`
}

type CConfigInfo struct {
	MqttLoadBalanceInfo CLoadBalanceInfo `json:"mqtt-loadbalance-info"`
}
