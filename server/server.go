package main

import (
	"bytes"
	"encoding/gob"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type FileServer struct {
	listenAddress string
}

func NewFileServer(listenAddr string) *FileServer {
	return &FileServer{
		listenAddress: listenAddr,
	}
}

func (fs *FileServer) start() {
	listener, err := net.Listen("tcp", fs.listenAddress)
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
	// var size int64
	// filenameBytes := make([]byte, 35)
	var filename string

	// filenameBuffer := new(bytes.Buffer)

	time.Sleep(5 * time.Second)
	// get read filename from connection
	// err := binary.Read(conn, binary.LittleEndian, &filename)

	err := gob.NewDecoder(conn).Decode(&filename)
	if err != nil {
		log.Panicln(err)
		return
	}

	// filename := string(filenameBytes)
	log.Println(filename)

	// return

	file, err := os.Create(filename)

	if err != nil {
		log.Panic(err)

		return
	}

	for {

		n, err := io.Copy(buf, conn)

		// if err == io.EOF {
		// 	log.Panicln(err)
		// 	break
		// }

		if err != nil {
			if err == io.EOF {
				log.Panicln(err)
				break
			}

			log.Panic(err)
			break
		}

		_, err = file.Write(buf.Bytes()[:n])
		log.Printf("error writing to file : %v\n", err)
	}

	log.Printf("File : %s  successfully received and saved.\n", filename)
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
