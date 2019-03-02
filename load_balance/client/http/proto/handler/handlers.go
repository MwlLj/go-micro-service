package handler

import (
	"encoding/json"
	"github.com/MwlLj/go-micro-service/load_balance/client/http/client"
	cfg "github.com/MwlLj/go-micro-service/load_balance/client/http/config"
	"github.com/MwlLj/go-micro-service/load_balance/client/http/proto/common"
	"github.com/MwlLj/go-micro-service/load_balance/client/http/proto/config"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func StartHandler(server common.IRun, w http.ResponseWriter, r *http.Request) {
	code := 0
	desc := "ok"
	for {
		if r.Method == http.MethodPost {
			request, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Printf("read request body error, err: %v\n", err)
				code = 100
				desc = "read request body error"
				break
			}
			startRequest := config.CRemoteProtoStartRequest{}
			err = json.Unmarshal([]byte(request), &startRequest)
			if err != nil {
				log.Printf("decode start request error, err: %v\n", err)
				code = 101
				desc = "decode start request error"
				break
			}
			cli := client.New(&startRequest.DispersedInfo)
			server.SetClient(cli)
		} else {
			log.Printf("request method error, method: %s\n", r.Method)
			code = 102
			desc = "method not allow"
			break
		}
		break
	}
	w.WriteHeader(http.StatusOK)
	response := config.CRemoteProtoStartResponse{}
	response.Error = code
	response.ErrorString = desc
	res, err := json.Marshal(&response)
	if err != nil {
		log.Printf("encoding response string error, err: %v\n", err)
		io.WriteString(w, "json encoding response error")
	}
	io.WriteString(w, string(res))
}

func GetConnectHandler(server common.IRun, w http.ResponseWriter, r *http.Request) {
	code := 0
	desc := "ok"
	var loadBalanceNetInfo *cfg.CLoadBalanceNetInfo
	var err error = nil
	for {
		if r.Method == http.MethodGet {
			cli := server.Client()
			if cli == nil {
				log.Println("start request never is send")
				code = 110
				desc = "please send start first"
				break
			}
			loadBalanceNetInfo, err = cli.GetConnect()
			if err != nil {
				log.Printf("get loadbalance error, err: %v\n", err)
				code = 111
				desc = "get loadbalance error"
			}
		} else {
			log.Printf("request method error, method: %s\n", r.Method)
			code = 112
			desc = "method not allow"
			break
		}
		break
	}
	w.WriteHeader(http.StatusOK)
	response := config.CRemoteProtoGetConnectResponse{}
	response.Error = code
	response.ErrorString = desc
	response.LoadBalanceNetInfo = *loadBalanceNetInfo
	res, err := json.Marshal(&response)
	if err != nil {
		log.Printf("encoding response string error, err: %v\n", err)
		io.WriteString(w, "json encoding response error")
	}
	io.WriteString(w, string(res))
}

func RegisterServiceHandler(server common.IRun, w http.ResponseWriter, r *http.Request) {
	code := 0
	desc := "ok"
	for {
		if r.Method == http.MethodPost {
			cli := server.Client()
			if cli == nil {
				log.Println("start request never is send")
				code = 120
				desc = "please send start first"
				break
			}
			err := cli.RegisterService()
			if err != nil {
				log.Printf("register service error, err: %v\n", err)
				code = 121
				desc = "register service error"
			}
		} else {
			log.Printf("request method error, method: %s\n", r.Method)
			code = 122
			desc = "method not allow"
			break
		}
		break
	}
	w.WriteHeader(http.StatusOK)
	response := config.CRemoteProtoRegisterServiceResponse{}
	response.Error = code
	response.ErrorString = desc
	res, err := json.Marshal(&response)
	if err != nil {
		log.Printf("encoding response string error, err: %v\n", err)
		io.WriteString(w, "json encoding response error")
	}
	io.WriteString(w, string(res))
}
