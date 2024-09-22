package main

import (
	"fmt"
	"net"
)

func handleConn(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Received: %s\n", buffer[:n])
	}
}

func runServer() {
	server, err := net.Listen("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(conn)
	}
}

func main() {
	runServer()
}
