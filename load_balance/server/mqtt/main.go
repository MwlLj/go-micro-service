package main

import (
	"./impl"
	"flag"
)

func main() {
	var configPath *string = flag.String("c", "mqtt-load-balance.cfg", "config file path")
	flag.Parse()

	server := impl.CServer{}
	server.Start(*configPath)
}
