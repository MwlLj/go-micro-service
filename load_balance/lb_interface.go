package load_balance

type CLoadBlance interface {
	GetMasterNode()
	RoundRobin()
	WeightRoundRobin()
	Random()
	WeightRandom()
	IpHash()
	UrlHash()
	LeastConnections()
}
