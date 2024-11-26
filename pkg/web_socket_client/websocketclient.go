package websocketclient

import (
	"encoding/binary"
	"net"
)

type Client struct {
	socket net.Conn
}

func NewClient(socket net.Conn) *Client {
	return &Client{
		socket: socket,
	}
}

func (c *Client) Read() ([]byte, error){
	finAndOpcode := make([]byte, 1)
	
	_, err := c.socket.Read(finAndOpcode)
	if err != nil {
		return nil, err
	}

	maskAndLengthIndicator := make([]byte, 1)
	_, err = c.socket.Read(maskAndLengthIndicator)
	if err != nil {
		return nil, err
	}

	lengthIndicator := maskAndLengthIndicator[0] & 0x7F
	
	var length int

	if lengthIndicator <= 125 {
		length = int(lengthIndicator)
	} else if lengthIndicator == 126 {
		lengthBytes := make([]byte, 2)
		_, err := c.socket.Read(lengthBytes)
		if err != nil {
			return nil, err
		}
		
		length = int(binary.BigEndian.Uint16(lengthBytes))
	} else if lengthIndicator == 127 {
		lengthBytes := make([]byte, 8)
		_, err := c.socket.Read(lengthBytes)
		if err != nil {
			return nil, err
		}

		length = int(binary.BigEndian.Uint64(lengthBytes))
	}
	
	maskKeys := make([]byte, 4)
	
	_, err = c.socket.Read(maskKeys)
	if err != nil {
		return nil, err
	}

	encoded := make([]byte, length)
	
	_, err = c.socket.Read(encoded)
	if err != nil {
		return nil, err
	}

	decoded := make([]byte, length)
	
	for i := 0; i < length; i++ {
		decoded[i] = encoded[i] ^ maskKeys[i%4]
	}

	return decoded, nil
}

func (c *Client) Send(message string) error {
	bytes := []byte{129}

	size := len(message)

	if size <= 125 {
		bytes = append(bytes, byte(size))
	} else if size < 1 << 16 {
		bytes = append(bytes, 126)

		sizeBytes := make([]byte, 2)

		binary.BigEndian.PutUint16(sizeBytes, uint16(size))

		bytes = append(bytes, sizeBytes...)
	} else {
		bytes = append(bytes, 127)

		sizeBytes := make([]byte, 8)
		
		binary.BigEndian.PutUint64(sizeBytes, uint64(size))

		bytes = append(bytes, sizeBytes...)
	}

	bytes = append(bytes, []byte(message)...)

	_, err := c.socket.Write(bytes)
	return err
}



