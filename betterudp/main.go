package main

import (
	"betterudp/client"
	"betterudp/server"

)

func main() {
	go server.Server(":1234")
	client.Client("127.0.0.1:1234")
}