package main

import (
	"log"
	"os"
)


func main() {
	fs := &FileServer{}
	
	go func() {
		
		file, err  := os.Open("The Rust Programming Language.pdf")
		if err != nil {
			log.Panicln(err)
		}
		
		
		serverClient, err := NewClient("tcp", ":8000")
		if err != nil {
			log.Panic(err)
		}
		
		// write file to server		
		err = serverClient.writeBytes(file)
		if err != nil {
			log.Fatal(err)
		}
	}()
	
	
	log.Println("Starting server")
	fs.start()
}
