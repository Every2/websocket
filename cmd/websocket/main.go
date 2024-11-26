package main

import (
	"log"
	"os"

	websocketserver "github.com/Every2/websocket/pkg/web_socket_server"
)

func main() {
	path := "/"
	port := 4567
	host := "localhost"
	server := websocketserver.NewServer(path, port, host)

	ln, err := server.Init()
	if err != nil {
		log.Fatalf("Server error: %v", err)
		os.Exit(1)
	}

	log.Printf("Servidor is working on %s:%d", host, port)
	for {
		server.Accept(ln)
	}
}
