package main

import (
	"flag"
	"github.com/MwlLj/go-micro-service/load_balance/client/http/proto/config"
	"github.com/MwlLj/go-micro-service/load_balance/client/http/proto/run"
	"log"
)

func main() {
	var configPath *string = flag.String("c", "http-loadbalance-client.cfg", "config file path")
	reader := config.CReader{}
	info, err := reader.Read(*configPath)
	if err != nil {
		log.Printf("load config error, err: %v\n", err)
		return
	}
	runner := run.CRun{}
	runner.Start(info)
}
