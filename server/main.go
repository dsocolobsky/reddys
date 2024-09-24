package main

import (
	"fmt"
	"github.com/dsocolobsky/reddys/internal"
	"net"
	"strings"
)

var database *internal.Database

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
	case "SET":
		if len(arr) == 3 {
			database.Set(arr[1], arr[2])
			message = internal.CraftSimpleString("OK")
		} else {
			message = internal.CraftSimpleError("wrong number of arguments for 'set' command")
		}
	case "GET":
		if len(arr) == 2 {
			value := database.Get(arr[1])
			if value == "" {
				message = internal.CraftNullString()
			} else {
				message = internal.CraftBulkString(value)
			}
		} else {
			message = internal.CraftSimpleError("wrong number of arguments for 'get' command")
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
	database = internal.NewDatabase()

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
