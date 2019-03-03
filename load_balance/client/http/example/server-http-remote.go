package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	bl "github.com/MwlLj/go-micro-service/load_balance"
	"github.com/MwlLj/go-micro-service/load_balance/client/http/config"
	cfg "github.com/MwlLj/go-micro-service/load_balance/client/http/proto/config"
	s "github.com/MwlLj/go-micro-service/service_discovery_nocache"
	"github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var _ = log.Println
var _ = fmt.Println

func AddPictureHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello too")
}

func main() {
	/*
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
	*/

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
	uid, err := uuid.NewV4()
	if err != nil {
		fmt.Println(err)
		return
	}
	serverUniqueCode := uid.String()
	serviceInfo := config.CServiceInfo{}
	serviceInfo.ServerName = "ress"
	serviceInfo.ServerVersion = "1.0"
	var serviceServiceDiscoveryConns []config.CServiceDiscoveryNet
	serviceServiceDiscoveryConn := config.CServiceDiscoveryNet{
		ServerHost: "localhost",
		ServerPort: 2182,
		ServiceId:  "server_1",
	}
	serviceServiceDiscoveryConns = append(serviceServiceDiscoveryConns, serviceServiceDiscoveryConn)
	serviceInfo.ServiceDiscoveryConns = serviceServiceDiscoveryConns
	serviceInfo.PathPrefix = "taobao-service"
	serviceInfo.ServerMode = s.ServerModeZookeeper
	serviceInfo.ServerUniqueCode = serverUniqueCode
	serviceInfo.Weight = 1
	serviceInfo.Host = "localhost"
	serviceInfo.Port = 20100
	serviceInfo.UserName = ""
	serviceInfo.UserPwd = ""
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
	registerServiceResponse, err := send("localhost", 15000, "/service/register", http.MethodPost, "")
	if err != nil {
		log.Fatalf("send /service/register error, err: %v\n", err)
	}
	registerServiceRes := cfg.CRemoteProtoRegisterServiceResponse{}
	err = json.Unmarshal([]byte(registerServiceResponse), &registerServiceRes)
	if err != nil {
		log.Fatalf("decoding json error, err: %v\n", err)
	}
	// start http server
	mux := http.NewServeMux()
	mux.HandleFunc("/picture", AddPictureHandler)
	addrBuffer := bytes.Buffer{}
	addrBuffer.WriteString(serviceInfo.Host)
	addrBuffer.WriteString(":")
	addrBuffer.WriteString(strconv.FormatInt(int64(serviceInfo.Port), 10))
	http.ListenAndServe(addrBuffer.String(), mux)
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
