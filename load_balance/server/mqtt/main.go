package main

import (
	"flag"
	"github.com/MwlLj/go-micro-service/load_balance/server/mqtt/impl"
)

func main() {
	var configPath *string = flag.String("c", "mqtt-load-balance.cfg", "config file path")
	flag.Parse()

	server := impl.CServer{}
	server.Start(*configPath)
}
