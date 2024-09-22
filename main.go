package main

import (
	"fmt"
	"net"
	"strings"
)

func simpleString(message string) string {
	return fmt.Sprintf("+%s\r\n", message)
}

func simpleError(message string) string {
	return fmt.Sprintf("-%s\r\n", message)
}

func bulkString(message string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(message), message)
}

func handleMessage(conn net.Conn, message string) {
	if message == "" {
		return
	}

	words := strings.SplitN(message, " ", 2)
	if len(words) == 0 {
		return
	}
	command := strings.ToUpper(strings.TrimSpace(words[0])) // Redis commands are case-insensitive
	fmt.Println("Command: ", command)
	switch command {
	case "PING":
		if len(words) > 1 {
			message = bulkString(words[1])
		} else {
			message = simpleString("PONG")
		}
	default:
		message = simpleError("ERR unknown command")
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

func runServer() {
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

func main() {
	runServer()
}
