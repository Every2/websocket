package websocketserver

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"

	websocketclient "github.com/Every2/websocket/pkg/web_socket_client"
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
	address := fmt.Sprintf("%s:%d", s.host, s.port)
	connection, err := net.Listen("tcp", address)

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

	if s.sendHandShake(connection) {
		websocketclient.NewClient(connection)
		return
	}
}

func (s *Server) sendHandShake(connection net.Conn) bool {
	reader := bufio.NewReader(connection)

	req, err := reader.ReadString('\n')

	if err != nil {
		return false
	}

	header, err := getHeader(reader)

	if err != nil {
		return false
	}

	req_rgx := regexp.MustCompile(`GET ` + regexp.QuoteMeta(s.path) + ` HTTP/1.1`)
	key_rgx := regexp.MustCompile(`Sec-WebSocket-Key: (.*)\r\n`)

	if req_rgx.MatchString(req) && key_rgx.MatchString(header) {
		matches := key_rgx.FindStringSubmatch(header)

		if len(matches) > 1 {
			new_accept := createAccept(matches[1])
			sendHandShakeResponse(connection, new_accept)
			return true
		}
	}

	connection.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
	return false
}

func getHeader(reader *bufio.Reader) (string, error) {
	var header strings.Builder

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			return "", err
		}

		if line == "\r\n" {
			break
		}

		header.WriteString(line)
	}

	return header.String(), nil
}

func createAccept(key string) string {

	magicString := "put something here"
	
	input := key + magicString

	hash := sha1.New()
	hash.Write([]byte(input))

	digest := hash.Sum(nil)

	accept := base64.StdEncoding.EncodeToString(digest)

	return accept
}

func sendHandShakeResponse(connection net.Conn, acceptString string) {
	response := fmt.Sprintf("HTTP/1.1 101 Switching Protocols\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Accept: %s\r\n", acceptString)

	connection.Write([]byte(response))
}
