package websocketserver

import (
	"log"
	"net"
)

type Server struct {
	path string
	port int
	host string
}

func NewServer(path string, port int, host string) *Server {
	return &Server{
		path: path,
		port: port,
		host: host,
	}
}

func (s *Server) Init() (*net.TCPListener, error) {
	connection, err := net.Listen("tcp", s.path+":"+string(s.port))

	if err != nil {
		return nil, err
	}

	return connection.(*net.TCPListener), nil
}

func (s *Server) Accept(listener *net.TCPListener) {
	connection, err := listener.Accept()

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	defer connection.Close()

	if sendHandShake(conn) {
		NewClient(conn)
		return
	}
}
