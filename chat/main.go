package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Galdoba/ffstuff/chat/server"
)

const (
	address = ":8888"
)

/*
func main2() {
	s := NewServer()
	go s.run()

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}
	defer listener.Close()
	log.Printf("started server on %s", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("unable to accept connection %s", err.Error())
			continue
		}

		go s.NewClient(conn)
	}

}
*/
func main() {

	s := server.NewServer()
	go s.Run()

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}
	defer listener.Close()
	log.Printf("started server on %s", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("unable to accept connection %s", err.Error())
			continue
		}

		go s.NewClient(conn)
		fmt.Println("END cycle")
	}

}

/*
curl -X POST -H "Content-Type: text/plain" --data "this is raw data" http://78.41.xx.xx:7778/
curl -X POST -H "Content-Type: text/plain" --data "test text" http://:8888
*/
