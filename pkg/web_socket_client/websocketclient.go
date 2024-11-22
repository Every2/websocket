package websocketclient

import "net"

type Client struct {
	socket net.Conn
}

func NewClient(socket net.Conn) *Client {
	return &Client{
		socket: socket,
	}
}


