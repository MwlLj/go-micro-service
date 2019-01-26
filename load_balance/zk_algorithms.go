package load_balance

import (
	proto "../common_proto"
	"errors"
	"fmt"
	"hash/crc32"
	"math/rand"
	"strconv"
	"strings"
	"sync"
)

var _ = fmt.Println
var _ = errors.New
var _ = strings.Join
var _ = strconv.ParseInt

type CRoundRobin struct {
	m_mutex      sync.Mutex
	m_loadBlance ILoadBlance
	m_nodeLength int
	m_nodeIndex  int
}

func (this *CRoundRobin) init(loadBlance ILoadBlance) error {
	this.m_loadBlance = loadBlance
	return nil
}

func (this *CRoundRobin) Get(serverName string, extraData interface{}) (*proto.CNodeData, error) {
	item, err := this.m_loadBlance.FindServerData(serverName)
	if err != nil {
		fmt.Println("[ERROR] server not find")
		return nil, err
	}
	length := len(*item.NormalNodes)
	if length == 0 {
		return nil, errors.New("normal node is not exist")
	}
	this.m_mutex.Lock()
	if length != this.m_nodeLength || this.m_nodeLength == 0 {
		this.m_nodeLength = length
		this.m_nodeIndex = 0
	}
	data := (*item.NormalNodes)[this.m_nodeIndex]
	this.m_nodeIndex += 1
	if this.m_nodeIndex > this.m_nodeLength-1 {
		this.m_nodeIndex = 0
	}
	this.m_mutex.Unlock()
	return &data, nil
}

type CWeightRoundRobin struct {
	m_loadBlance ILoadBlance
	m_mutex      sync.Mutex
	m_nodeLength int
	m_nodeIndex  int
	m_curWeight  int
}

func (this *CWeightRoundRobin) init(loadBlance ILoadBlance) error {
	this.m_loadBlance = loadBlance
	return nil
}

func (this *CWeightRoundRobin) Get(serverName string, extraData interface{}) (*proto.CNodeData, error) {
	item, err := this.m_loadBlance.FindServerData(serverName)
	if err != nil {
		fmt.Println("[ERROR] server not find")
		return nil, err
	}
	this.m_mutex.Lock()
	length := len(*item.NormalNodes)
	if length != this.m_nodeLength || this.m_nodeLength == 0 {
		this.m_nodeLength = length
		this.m_nodeIndex = 0
		this.m_curWeight = 0
	}
	data := (*item.NormalNodes)[this.m_nodeIndex]
	weight := data.Weight
	if this.m_curWeight == 0 && weight > 0 {
		this.m_curWeight = weight
	}
	if weight > 0 {
		this.m_curWeight -= 1
	}
	if this.m_curWeight <= 0 {
		this.m_nodeIndex += 1
	}
	if this.m_nodeIndex > this.m_nodeLength-1 {
		this.m_nodeIndex = 0
	}
	this.m_mutex.Unlock()
	return &data, nil
}

func randomInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

type CRandom struct {
	m_loadBlance ILoadBlance
	m_nodeLength int
}

func (this *CRandom) init(loadBlance ILoadBlance) error {
	this.m_loadBlance = loadBlance
	return nil
}

func (this *CRandom) Get(serverName string, extraData interface{}) (*proto.CNodeData, error) {
	item, err := this.m_loadBlance.FindServerData(serverName)
	if err != nil {
		fmt.Println("[ERROR] server not find")
		return nil, err
	}
	length := len(*item.NormalNodes)
	if length == 0 {
		return nil, errors.New("[ERROR] normal node is null")
	}
	if length != this.m_nodeLength || this.m_nodeLength == 0 {
		this.m_nodeLength = length
	}
	randomValue := randomInt(0, this.m_nodeLength)
	data := (*item.NormalNodes)[randomValue]
	return &data, nil
}

type CWeightRandom struct {
	m_loadBlance ILoadBlance
	m_nodeLength int
}

func (this *CWeightRandom) init(loadBlance ILoadBlance) error {
	this.m_loadBlance = loadBlance
	return nil
}

func (this *CWeightRandom) Get(serverName string, extraData interface{}) (*proto.CNodeData, error) {
	item, err := this.m_loadBlance.FindServerData(serverName)
	if err != nil {
		fmt.Println("[ERROR] server not find")
		return nil, err
	}
	length := len(*item.NormalNodes)
	if length == 0 {
		return nil, errors.New("[ERROR] normal node is null")
	}
	if length != this.m_nodeLength || this.m_nodeLength == 0 {
		this.m_nodeLength = length
	}
	nodes := this.rebuildNormalNodes(item.NormalNodes)
	rebuildLength := len(*nodes)
	randomValue := randomInt(0, rebuildLength)
	data := (*nodes)[randomValue]
	return &data, nil
}

func (*CWeightRandom) rebuildNormalNodes(normals *[]proto.CNodeData) *[]proto.CNodeData {
	var nodes []proto.CNodeData
	for _, normal := range *normals {
		weight := normal.Weight
		for i := 0; i < weight; i++ {
			nodes = append(nodes, normal)
		}
	}
	return &nodes
}

func toHash(b []byte) int {
	v := int(crc32.ChecksumIEEE(b))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

type CIpHash struct {
	m_loadBlance ILoadBlance
}

func (this *CIpHash) init(loadBlance ILoadBlance) error {
	this.m_loadBlance = loadBlance
	return nil
}

func (this *CIpHash) Get(serverName string, extraData interface{}) (*proto.CNodeData, error) {
	if extraData == nil {
		return nil, errors.New("ip is null, you should give extraData ip (type: string)")
	}
	item, err := this.m_loadBlance.FindServerData(serverName)
	if err != nil {
		fmt.Println("[ERROR] server not find")
		return nil, err
	}
	length := len(*item.NormalNodes)
	if length == 0 {
		return nil, errors.New("normal node is null")
	}
	hash := toHash([]byte(extraData.(string)))
	index := hash % length
	data := (*item.NormalNodes)[index]
	return &data, nil
}

type CUrlHash struct {
	m_loadBlance ILoadBlance
}

func (this *CUrlHash) init(loadBlance ILoadBlance) error {
	this.m_loadBlance = loadBlance
	return nil
}

func (this *CUrlHash) Get(serverName string, extraData interface{}) (*proto.CNodeData, error) {
	if extraData == nil {
		return nil, errors.New("url is null, you should give extraData url (type: string)")
	}
	item, err := this.m_loadBlance.FindServerData(serverName)
	if err != nil {
		fmt.Println("[ERROR] server not find")
		return nil, err
	}
	length := len(*item.NormalNodes)
	if length == 0 {
		return nil, errors.New("normal node is null")
	}
	hash := toHash([]byte(extraData.(string)))
	index := hash % length
	data := (*item.NormalNodes)[index]
	return &data, nil
}

type CLeastConnections struct {
	m_loadBlance       ILoadBlance
	m_algorithm        INormalNodeAlgorithm
	m_normalNodeRecord sync.Map
}

type leastConnectionItem struct {
	data  *proto.CNodeData
	times int
}

func (this *CLeastConnections) init(loadBlance ILoadBlance) error {
	this.m_loadBlance = loadBlance
	this.m_loadBlance.SetCallback(this, nil)
	return nil
}

func (this *CLeastConnections) Get(serverName string, extraData interface{}) (*proto.CNodeData, error) {
	item, err := this.m_loadBlance.FindServerData(serverName)
	if err != nil {
		fmt.Println("[ERROR] server not find")
		return nil, err
	}
	length := len(*item.NormalNodes)
	if length == 0 {
		return nil, errors.New("normal node is null")
	}
	v, ok := this.m_normalNodeRecord.Load(serverName)
	var items *[]leastConnectionItem = nil
	if ok {
		items = v.(*[]leastConnectionItem)
	} else {
		items = new([]leastConnectionItem)
	}
	if item.IsChanged == true || !ok {
		fmt.Println("[INFO] items changed -> delete first, then add")
		*items = (*items)[0:0]
		// changed or not exist -> delete first, then add
		this.m_normalNodeRecord.Delete(serverName)
		for _, node := range *item.NormalNodes {
			tmp := node
			*items = append(*items, leastConnectionItem{data: &tmp, times: 0})
		}
		this.m_normalNodeRecord.Store(serverName, items)
	}
	firstItem := (*items)[0]
	min := firstItem.times
	minValue := firstItem.data
	for _, it := range *items {
		if it.times < min {
			min = it.times
			minValue = it.data
		}
	}
	fmt.Println("[INFO] mintimes: ", min)
	// update times
	for i, it := range *items {
		itHash := this.m_loadBlance.nodeData2hash(it.data)
		minHash := this.m_loadBlance.nodeData2hash(minValue)
		if itHash == minHash {
			(*items)[i].times += 1
			break
		}
	}
	this.m_normalNodeRecord.Store(serverName, items)
	return minValue, nil
}

func (*CLeastConnections) MasterNodeChange(data *proto.CNodeData, userData interface{}) {
}

func (*CLeastConnections) NormalNodeChange(data *proto.CNodeData, userData interface{}) {
}

func (this *CLeastConnections) ServerBeDeleted(serverName *string, userData interface{}) {
	this.m_normalNodeRecord.Delete(*serverName)
}
