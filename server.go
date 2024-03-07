package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

type FileServer struct{}

func (fs *FileServer) start() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	
	log.Println("Server started & Listening....")
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Panic(err)
		}

		go fs.readLoop(conn)

	}
}

func (fs *FileServer) readLoop(conn net.Conn) {
	buf := new(bytes.Buffer)
	for {
		var size int64
		binary.Read(conn, binary.LittleEndian, &size)
		_ , err := io.Copy(buf, conn)
		if err != nil {
			log.Panic(err)
		}

		fmt.Println(string(buf.Bytes()))

	}
}

//func sendFile(size int) error {
//	file := make([]byte, size)
//
//	_, err := io.ReadFull(rand.Reader, file)
//	if err != nil {
//		return err
//	}
//
//	conn, err := net.Dial("tcp", ":3000")
//	if err != nil {
//		return err
//	}
//
//	binary.Write(conn, binary.LittleEndian, int64(size))
//	n, err := io.CopyN(conn, bytes.NewBuffer(file), int64(size))
//	if err != nil {
//		return err
//	}
//
//	fmt.Printf("Written %d bytes to server\n", n)
//	return nil
//}