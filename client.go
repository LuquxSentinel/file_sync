package main

import (
	"encoding/binary"
	"io"
	"net"
	"os"
)

type client struct {
	connection net.Conn
}

func NewClient(network string,address string) (*client, error) {
	
//	dial network connection to [address] using network ["tcp" | ""] 
	connection, err := net.Dial(network, address)
	if err != nil {
//		dial error
		return nil, err
	}
	
	return &client{
		connection: connection,
	}, nil
}

func (c *client) close() error {
	return c.connection.Close()
}


//func writeFileBytes(file *os.File) {
//	
//	buf := make([]byte, 1024)
//	
//	for {
//		n, err := file.Read(buf)
//		if err ==io.EOF {
//			break
//		}
//		
//		_, err := 
//	}
//	
//	io.Reader().Read(buf)
//}

func (c *client) writeBytes(file *os.File) error {
	
//	file information [stats]
	fileStats, _ := file.Stat()
	
	
	sizeBuffer := make([]byte, 8)
	
	binary.LittleEndian.AppendUint16(sizeBuffer, uint16(fileStats.Size()))
	
	_, err := c.connection.Write(sizeBuffer)
	if err != nil {
		return err
	}
	
	
//	write file to connection [network ("tcp") server]
	_, err = io.Copy(c.connection, file)
	return err
}
