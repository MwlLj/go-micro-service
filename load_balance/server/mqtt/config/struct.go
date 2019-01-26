package config

type CBrokerNet struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type CBrokerNetInfo struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UserName string `json:"username"`
	UserPwd  string `json:"userpwd"`
}

type CServiceDiscoryNet struct {
	ServerHost string
	ServerPort int
	ServiceId  string
}

type CServiceDiscory struct {
	PathPrefix string               `json:"pathprefix"`
	ServerMode string               `json:"servermode"`
	ServerName string               `json:"servername"`
	NodeData   CBrokerNet           `json:"nodedata"`
	Weight     int                  `json:"weight"`
	Conns      []CServiceDiscoryNet `json:"conns"`
}

type CLoadBalanceInfo struct {
	PathPrefix          string               `json:"pathprefix"`
	ServerMode          string               `json:"servermode"`
	NormalNodeAlgorithm string               `json:"normalNodeAlgorithm"`
	Conns               []CServiceDiscoryNet `json:"conns"`
}

type CRecvBrokerInfo struct {
	Nets []CBrokerNetInfo `json:"nets"`
}

type CRouterRule struct {
	Rule       string `json:"rule"`
	ServerName string `json:servername`
	IsMaster   bool   `json:"ismaster"`
}

type CRouterRuleInfo struct {
	Rules []CRouterRule `json:"rules"`
}

type CConfigInfo struct {
	BrokerRegisterInfo  CServiceDiscory  `json:"brokerreg"`
	ServiceRegisterInfo CLoadBalanceInfo `json:"servicereg"`
	RecvBrokerInfo      CRecvBrokerInfo  `json:"recvbrokerinfo"`
	RouterRuleInfo      CRouterRuleInfo  `json:"routerrule"`
}
