package config

import (
	cfg "github.com/MwlLj/go-micro-service/load_balance/client/http/config"
)

type CRemoteListenNetInfo struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type CRemoteProtoStartProto struct {
	Url string `json:"url"`
}

type CRemoteProtoStartRequest struct {
	DispersedInfo cfg.CConfigInfo `json:"dispersed-info"`
}

type CRemoteProtoStartResponse struct {
	Error       int    `json:"error"`
	ErrorString string `json:"error-string"`
}

type CRemoteProtoGetConnectProto struct {
	Url string `json:"url"`
}

type CRemoteProtoGetConnectResponse struct {
	LoadBalanceNetInfo cfg.CLoadBalanceNetInfo `json:"load-balance-info"`
	Error              int                     `json:"error"`
	ErrorString        string                  `json:"error-string"`
}

type CRemoteProtoRegisterServiceProto struct {
	Url string `json:"url"`
}

type CRemoteProtoRegisterServiceRequest struct {
}

type CRemoteProtoRegisterServiceResponse struct {
	Error       int    `json:"error"`
	ErrorString string `json:"error-string"`
}

type CRemoteProtoConfigInfo struct {
	Net                 CRemoteListenNetInfo `json:"net"`
	StartRule           CRemoteProtoStartProto
	GetConnectRule      CRemoteProtoGetConnectProto      `json:"get-connect-proto"`
	RegisterServiceRule CRemoteProtoRegisterServiceProto `json:"register-service-proto"`
}
