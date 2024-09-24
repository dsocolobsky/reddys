package main

import (
	"fmt"
	"github.com/dsocolobsky/reddys/internal"
	"net"
	"strings"
)

func handleMessage(conn net.Conn, message string) {
	if message == "" {
		return
	}

	arr := internal.ReadArray(message)
	if len(arr) == 0 {
		fmt.Println("Empty array of commands!")
	}

	command := strings.ToUpper(strings.TrimSpace(arr[0]))

	switch command {
	case "PING":
		if len(arr) == 1 {
			message = internal.CraftSimpleString("PONG")
		} else if len(arr) == 2 {
			message = internal.CraftBulkString(arr[1])
		} else {
			message = internal.CraftSimpleError("wrong number of arguments for 'ping' command")
		}
	default:
		message = internal.CraftSimpleError("ERR unknown command")
	}

	_, err := conn.Write([]byte(message))
	if err != nil {
		panic(err)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			panic(err)
		}
		handleMessage(conn, string(buffer[:n]))
		fmt.Printf("Received: %s\n", buffer[:n])
	}
}

func main() {
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
