package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type client struct {
	connection net.Conn
	watcher    *fsnotify.Watcher
}

func NewClient(network string, address string) (*client, error) {

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

func (c *client) writeBytes(filename string) error {

	log.Println("----")
	log.Println(filename)

	file, err := os.Open(filename)
	if err != nil {
		log.Panicln(err)
	}

	//	file information [stats]
	fileStats, _ := file.Stat()

	// sizeBuffer := make([]byte, 8)
	binary.Write(c.connection, binary.LittleEndian, fileStats.Name())

	// binary.LittleEndian.AppendUint16(sizeBuffer, uint16(fileStats.Size()))

	// _, err = c.connection.Write(sizeBuffer)
	// if err != nil {
	// 	return err
	// }

	//	write file to connection [network ("tcp") server]
	_, err = io.Copy(c.connection, file)
	return err
}

func (c *client) fileChangeListener() {
	log.Println("Listening to  changes in directory")
	//	file watcher
	c.watcher, _ = fsnotify.NewWatcher()
	defer c.watcher.Close()

	baseDir, err := os.UserHomeDir()
	if err != nil {
		log.Panicln(err)
	}

	log.Println(baseDir)
	err = filepath.Walk(fmt.Sprintf("%s/Documents/sync/", baseDir), c.watchDir)
	if err != nil {
		log.Panicln(err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {

			// watch for events
			case event := <-c.watcher.Events:
				// if watcher emits event
				// write changes to tcp server
				log.Println("jkfhjkdsfhjkdshfjkdshjkh")
				c.writeBytes(event.Name)

			case err := <-c.watcher.Errors:
				// if watcher emits error
				log.Println(err)
			}
		}
	}()

	<-done
}

func (c *client) watchDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return c.watcher.Add(path)
	}

	return nil
}
