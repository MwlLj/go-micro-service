package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	bl "github.com/MwlLj/go-micro-service/load_balance"
	"github.com/MwlLj/go-micro-service/load_balance/client/http/config"
	cfg "github.com/MwlLj/go-micro-service/load_balance/client/http/proto/config"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var _ = log.Println
var _ = fmt.Println

func main() {
	info := config.CConfigInfo{}
	loadBalanceInfo := config.CLoadBalanceInfo{}
	loadBalanceInfo.PathPrefix = "micro-service"
	loadBalanceInfo.ServerMode = bl.ServerModeZookeeper
	loadBalanceInfo.NormalNodeAlgorithm = bl.AlgorithmRoundRobin
	loadBalanceInfo.HttpLoadBalanceServerName = "http-nginx"
	var conns []config.CServiceDiscoveryNet
	conn := config.CServiceDiscoveryNet{
		ServerHost: "localhost",
		ServerPort: 2182,
		ServiceId:  "server_1",
	}
	conns = append(conns, conn)
	loadBalanceInfo.Conns = conns
	info.HttpLoadBalanceInfo = loadBalanceInfo
	// service init
	serviceInfo := config.CServiceInfo{}
	serviceInfo.ServerName = "cgws"
	serviceInfo.ServerVersion = "1.0"
	info.ServiceInfo = serviceInfo
	// start init
	startReq := cfg.CRemoteProtoStartRequest{}
	startReq.DispersedInfo = info
	startRequest, err := json.Marshal(&startReq)
	if err != nil {
		log.Fatalf("encoding error, err: %v\n", err)
	}
	_, err = send("localhost", 15000, "/server/init", http.MethodPost, string(startRequest))
	if err != nil {
		log.Fatalf("send server/init error, err: %v\n", err)
	}
	// get connect
	getConnectResponse, err := send("localhost", 15000, "/connect/loadbalance", http.MethodGet, "")
	if err != nil {
		log.Fatalf("send /connect/loadbalance error, err: %v\n", err)
	}
	getConnectRes := cfg.CRemoteProtoGetConnectResponse{}
	err = json.Unmarshal([]byte(getConnectResponse), &getConnectRes)
	if err != nil {
		log.Fatalf("decoding json error, err: %v\n", err)
	}
	loadBalanceNetInfo := getConnectRes.LoadBalanceNetInfo
	// send
	resBody, err := send("localhost", loadBalanceNetInfo.Port, "/picture", http.MethodPost, "hello")
	if err != nil {
		log.Fatalf("read resBody error, err: %v", err)
	}
	fmt.Printf("recv response, response: %s\n", resBody)
}

func send(host string, port int, url string, method string, request string) (string, error) {
	addrBuffer := bytes.Buffer{}
	addrBuffer.WriteString("http://")
	addrBuffer.WriteString(host)
	addrBuffer.WriteString(":")
	addrBuffer.WriteString(strconv.FormatInt(int64(port), 10))
	addrBuffer.WriteString(url)
	body := strings.NewReader(request)
	req, err := http.NewRequest(method, addrBuffer.String(), body)
	if err != nil {
		log.Fatalf("new request error, err: %v\n", err)
	}
	cli := http.DefaultClient
	res, err := cli.Do(req)
	if err != nil {
		log.Fatalf("send http request error, err: %v\n", err)
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("read response error, err: %v\n", err)
	}
	return string(response), nil
}
