package run

import (
	"bytes"
	"fmt"
	"github.com/MwlLj/go-micro-service/load_balance/client/http/client"
	"github.com/MwlLj/go-micro-service/load_balance/client/http/proto/config"
	"github.com/MwlLj/go-micro-service/load_balance/client/http/proto/handler"
	"net/http"
	"strconv"
)

var _ = fmt.Print

type CRun struct {
	m_client client.IClient
}

func (this *CRun) SetClient(cli client.IClient) {
	this.m_client = cli
}

func (this *CRun) Client() client.IClient {
	return this.m_client
}

func (this *CRun) Start(configInfo *config.CRemoteProtoConfigInfo) {
	// register handlers
	mux := http.NewServeMux()
	mux.HandleFunc(configInfo.StartRule.Url, func(w http.ResponseWriter, r *http.Request) {
		handler.StartHandler(this, w, r)
	})
	mux.HandleFunc(configInfo.GetConnectRule.Url, func(w http.ResponseWriter, r *http.Request) {
		handler.GetConnectHandler(this, w, r)
	})
	mux.HandleFunc(configInfo.RegisterServiceRule.Url, func(w http.ResponseWriter, r *http.Request) {
		handler.RegisterServiceHandler(this, w, r)
	})
	// start server
	addrBuffer := bytes.Buffer{}
	addrBuffer.WriteString(configInfo.Net.Host)
	addrBuffer.WriteString(":")
	addrBuffer.WriteString(strconv.FormatInt(int64(configInfo.Net.Port), 10))
	http.ListenAndServe(addrBuffer.String(), mux)
}
