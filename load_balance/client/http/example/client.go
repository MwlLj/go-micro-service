package main

import (
	"bytes"
	"fmt"
	bl "github.com/MwlLj/go-micro-service/load_balance"
	"github.com/MwlLj/go-micro-service/load_balance/client/http/client"
	"github.com/MwlLj/go-micro-service/load_balance/client/http/config"
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
	cli := client.New(&info)
	loadBalanceNetInfo, err := cli.GetConnect()
	if err != nil {
		log.Fatalln(err)
		return
	}
	// send
	urlBuffer := bytes.Buffer{}
	urlBuffer.WriteString("http://")
	urlBuffer.WriteString(loadBalanceNetInfo.Host)
	urlBuffer.WriteString(":")
	urlBuffer.WriteString(strconv.FormatInt(int64(loadBalanceNetInfo.Port), 10))
	urlBuffer.WriteString("/picture")
	body := strings.NewReader("hello")
	req, err := http.NewRequest(http.MethodPost, urlBuffer.String(), body)
	if err != nil {
		log.Fatalf("new request error, err: %v", err)
	}
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("send http request error, err: %v", err)
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("read resBody error, err: %v", err)
	}
	fmt.Printf("recv response, response: %s\n", resBody)
}
