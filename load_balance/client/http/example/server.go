package main

import (
	"bytes"
	"fmt"
	bl "github.com/MwlLj/go-micro-service/load_balance"
	"github.com/MwlLj/go-micro-service/load_balance/client/http/client"
	"github.com/MwlLj/go-micro-service/load_balance/client/http/config"
	s "github.com/MwlLj/go-micro-service/service_discovery_nocache"
	"github.com/satori/go.uuid"
	"io"
	"log"
	"net/http"
	"strconv"
)

var _ = fmt.Println

type CServer struct {
}

func CAddPictureHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello too")
}

func (this *CServer) Start() {
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
	cli := client.New(&info)
	err = cli.RegisterService()
	if err != nil {
		log.Fatalf("register service to discovery error, err: %v\n", err)
	}
	// start http server
	mux := http.NewServeMux()
	mux.HandleFunc("/picture", CAddPictureHandler)
	addrBuffer := bytes.Buffer{}
	addrBuffer.WriteString(serviceInfo.Host)
	addrBuffer.WriteString(":")
	addrBuffer.WriteString(strconv.FormatInt(int64(serviceInfo.Port), 10))
	http.ListenAndServe(addrBuffer.String(), mux)
}

func main() {
	server := CServer{}
	server.Start()
}
