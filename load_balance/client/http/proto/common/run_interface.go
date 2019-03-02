package common

import (
	"github.com/MwlLj/go-micro-service/load_balance/client/http/client"
)

type IRun interface {
	SetClient(cli client.IClient)
	Client() client.IClient
}
