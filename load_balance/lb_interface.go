package load_balance

type CLoadBlance interface {
	GetMasterNode(serverName string) *string
	RoundRobin(serverName string) *string
	WeightRoundRobin(serverName string) *string
	Random(serverName string) *string
	WeightRandom(serverName string) *string
	IpHash(serverName string) *string
	UrlHash(serverName string) *string
	LeastConnections(serverName string) *string
}
