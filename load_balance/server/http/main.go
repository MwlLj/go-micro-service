package main

import (
	"flag"
	"github.com/MwlLj/go-micro-service/load_balance/server/http/impl"
)

func main() {
	var configPath *string = flag.String("c", "http-load-balance.cfg", "config file path")
	flag.Parse()

	server := impl.CServer{}
	server.Start(*configPath)
}
