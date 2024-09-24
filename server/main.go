package main

import (
	"fmt"
	"github.com/dsocolobsky/reddys/internal"
	"net"
)

var handler *Handler

func toCommandArray(message string) []string {
	arr := internal.ReadArray(message)
	if len(arr) == 0 {
		fmt.Println("Empty array of commands!")
	}
	return arr
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			panic(err)
		}
		msg := string(buffer[:n])
		if msg == "" {
			continue
		}
		arr := toCommandArray(msg)
		res := handler.HandleCommand(arr)
		_, err = conn.Write([]byte(res))
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	handler = NewHandler(internal.NewMapDatabase())

	server, err := net.Listen("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}
	defer server.Close()

	fmt.Println("Server running on localhost:6379")
	for {
		conn, err := server.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(conn)
	}
}
