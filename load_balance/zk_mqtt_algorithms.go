package load_balance

import (
	// proto "../common_proto"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var _ = fmt.Println
var _ = errors.New
var _ = strings.Join
var _ = strconv.ParseInt
var _ = sync.NewCond

type CMqttRoundRobin struct {
	CRoundRobin
}

type CMqttWeightRoundRobin struct {
	CWeightRoundRobin
}

type CMqttRandom struct {
	CRandom
}

type CMqttWeightRandom struct {
	CWeightRandom
}

type CMqttIpHash struct {
	CIpHash
}

type CMqttUrlHash struct {
	CUrlHash
}

type CMqttLeastConnections struct {
	CLeastConnections
}
