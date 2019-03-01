package config

type CBrokerNetInfo struct {
	Host             string `json:"host"`
	Port             int    `json:"port"`
	UserName         string `json:"user-name"`
	UserPwd          string `json:"user-pwd"`
	ServerUniqueCode string `json:"server-uniquecode"`
}

type CServiceDiscoryNet struct {
	ServerHost string `json:"server-host"`
	ServerPort int    `json:"server-port"`
	ServiceId  string `json:"service-id"`
}

type CServiceDiscory struct {
	PathPrefix string               `json:"path-prefix"`
	ServerMode string               `json:"server-mode"`
	ServerName string               `json:"server-name"`
	NodeData   CBrokerNetInfo       `json:"nodedata"`
	Weight     int                  `json:"weight"`
	Conns      []CServiceDiscoryNet `json:"conns"`
}

type CLoadBalanceInfo struct {
	PathPrefix          string               `json:"path-prefix"`
	ServerMode          string               `json:"server-mode"`
	NormalNodeAlgorithm string               `json:"normal-node-algorithm"`
	Conns               []CServiceDiscoryNet `json:"conns"`
}

type CRecvBrokerInfo struct {
	Nets []CBrokerNetInfo `json:"nets"`
}

type CRouterRule struct {
	Rule       string `json:"rule"`
	ServerName string `json:"server-name"`
	IsMaster   bool   `json:"is-master"`
}

type CRouterRuleInfo struct {
	Rules []CRouterRule `json:"rules"`
}

type CConfigInfo struct {
	LoadBalanceInfo     CServiceDiscory  `json:"loadbalance-info"`
	ServiceRegisterInfo CLoadBalanceInfo `json:"service-reg"`
	RouterRuleInfo      CRouterRuleInfo  `json:"router-rule"`
}
