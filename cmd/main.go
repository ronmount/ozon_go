package main

import (
	"github.com/ronmount/ozon_go/cmd/backend_server"
)

func main() {
	server := backend_server.NewServer()
	server.Run()
}
