package config

type CLoadBalanceNetInfo struct {
	Host             string `json:"host"`
	Port             int    `json:"port"`
	UserName         string `json:"user-name"`
	UserPwd          string `json:"user-pwd"`
	ServerUniqueCode string `json:"server-uniquecode"`
}

type CServiceDiscoveryNet struct {
	ServerHost string `json:"server-host"`
	ServerPort int    `json:"server-port"`
	ServiceId  string `json:"service-id"`
}

type CLoadBalanceInfo struct {
	PathPrefix                string                 `json:"path-prefix"`
	ServerMode                string                 `json:"server-mode"`
	NormalNodeAlgorithm       string                 `json:"normal-node-algorithm"`
	HttpLoadBalanceServerName string                 `json:"mqtt-loadbalance-servername"`
	Conns                     []CServiceDiscoveryNet `json:"conns"`
}

type CServiceInfo struct {
	PathPrefix            string                 `json:"path-prefix"`
	ServerMode            string                 `json:"server-mode"`
	ServerName            string                 `json:"server-name"`
	ServerVersion         string                 `json:"server-version"`
	ServerUniqueCode      string                 `json:"server-uniquecode"`
	Weight                int                    `json:"weight"`
	Host                  string                 `json:"host"`
	Port                  int                    `json:"port"`
	UserName              string                 `json:"username"`
	UserPwd               string                 `json:"userpwd"`
	ServiceDiscoveryConns []CServiceDiscoveryNet `json:"service-discovery-conns"`
}

type CConfigInfo struct {
	HttpLoadBalanceInfo CLoadBalanceInfo `json:"http-loadbalance-info"`
	ServiceInfo         CServiceInfo     `json:"service-info"`
}
