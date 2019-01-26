package main

import (
	"./impl"
)

func main() {
	server := impl.CServer{}
	server.Start("mqtt-load-balance.cfg")
}
