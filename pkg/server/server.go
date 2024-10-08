package server

import (
	"fmt"
	"github.com/dsocolobsky/reddys/pkg/resp"
	"net"
	"strings"
)

type Server struct {
	handler   *Handler
	tcpserver net.Listener
}

func NewServer(handler *Handler) *Server {
	tcpserver, err := net.Listen("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}
	return &Server{
		handler:   handler,
		tcpserver: tcpserver,
	}
}

func (srv *Server) Close() {
	fmt.Println("Closing server...")
	err := srv.tcpserver.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println("Server closed")
}

func (srv *Server) Serve() {
	srv.handler.readFromDisk()

	fmt.Println("Server running on localhost:6379")
	for {
		conn, err := srv.tcpserver.Accept()
		if err != nil {
			panic(err)
		}
		go srv.handleConn(conn)
	}
}

func (srv *Server) handleConn(conn net.Conn) {
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
		fmt.Println(msg)
		arr, _ := resp.UnmarshalArray(msg)
		if len(arr) == 0 {
			fmt.Println("Empty array of commands!")
		}
		if len(arr) > 0 && writeCommands[strings.ToUpper(arr[0])] {
			srv.handler.Persist(msg)
		}
		res := srv.handler.HandleCommand(arr)
		_, err = conn.Write([]byte(res))
		if err != nil {
			panic(err)
		}
	}
}
