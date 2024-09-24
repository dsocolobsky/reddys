package main

import (
	"bufio"
	"github.com/dsocolobsky/reddys/internal"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	println("Connected to localhost:6379")

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		respMessage := internal.CraftBulkString(line)
		// Print resp message
		println(respMessage)
		// Send the message
		_, err = conn.Write([]byte(respMessage))
		if err != nil {
			panic(err)
		}
	}
}
