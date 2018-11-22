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

type CHttpRoundRobin struct {
	CRoundRobin
}

type CHttpWeightRoundRobin struct {
	CWeightRoundRobin
}

type CHttpRandom struct {
	CRandom
}

type CHttpWeightRandom struct {
	CWeightRandom
}

type CHttpIpHash struct {
	CIpHash
}

type CHttpUrlHash struct {
	CUrlHash
}

type CHttpLeastConnections struct {
	CLeastConnections
}
