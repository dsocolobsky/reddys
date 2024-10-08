package server

import (
	"fmt"
	"github.com/dsocolobsky/reddys/pkg/database"
	"github.com/dsocolobsky/reddys/pkg/resp"
	"net"
	"strings"
	"testing"
	"time"
)

// Utility function to send and receive commands
func sendCommand(conn *net.Conn, command string) (string, error) {
	words := strings.Split(command, " ")
	if len(words) == 0 {
		return "", fmt.Errorf("Empty command")
	}
	msg := resp.MarshalArrayOfBulkStrings(words)
	_, err := fmt.Fprintf(*conn, msg)
	if err != nil {
		return "", err
	}
	buffer := make([]byte, 1024)
	n, err := (*conn).Read(buffer)
	if err != nil {
		return "", err
	}
	res := string(buffer[:n])
	return res, nil
}

func createServer() *Server {
	handler := NewHandler(database.NewMapDatabase(), nil)
	server := NewServer(handler)
	return server
}

func assertCommandString(t *testing.T, conn *net.Conn, command string, expected string) {
	reply, err := sendCommand(conn, command)
	if err != nil {
		t.Fatalf("Failed to send %s command: %v", command, err)
	}
	res, _ := resp.UnmarshalRESP(reply)
	if res != expected {
		t.Fatalf("Expected %s, got %s", expected, reply)
	}
}

func assertCommandArrayResponse(t *testing.T, conn *net.Conn, command string, expected []string) {
	reply, err := sendCommand(conn, command)
	if err != nil {
		t.Fatalf("Failed to send %s command: %v", command, err)
	}
	res, _ := resp.UnmarshalArray(reply)
	if len(res) != len(expected) {
		t.Fatalf("Expected %d elements, got %d", len(expected), len(res))
	}
	for i, v := range expected {
		if v != res[i] {
			t.Fatalf("Expected %s, got %s", v, res[i])
		}
	}
}

func TestServer_Simple(t *testing.T) {
	server := createServer()
	go server.Serve()
	//defer server.Close() (check this later)
	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	assertCommandString(t, &conn, "PING", "PONG")
	assertCommandString(t, &conn, "SET foo bar", "OK")
	assertCommandString(t, &conn, "GET foo", "bar")
	assertCommandString(t, &conn, "GET bar", "")
	assertCommandString(t, &conn, "INCR val", "1")
	assertCommandString(t, &conn, "INCR val", "2")
	assertCommandString(t, &conn, "HSET player life 100", "OK")
	assertCommandString(t, &conn, "HSET player name duke", "OK")
	assertCommandString(t, &conn, "HGET player life", "100")
	assertCommandString(t, &conn, "HGET player name", "duke")
	assertCommandArrayResponse(t, &conn, "HGETALL player", []string{"life", "100", "name", "duke"})
	assertCommandString(t, &conn, "DBSIZE", "3")
	// GETDEL
	assertCommandString(t, &conn, "GETDEL val", "2")
	assertCommandString(t, &conn, "DBSIZE", "3")
	// GETSET
	assertCommandString(t, &conn, "GET foo", "bar")
	assertCommandString(t, &conn, "GETSET foo baz", "bar")
	assertCommandString(t, &conn, "GET foo", "baz")
}

func TestServer_TwoConnections(t *testing.T) {
	server := createServer()
	go server.Serve()
	time.Sleep(100 * time.Millisecond)

	conn1, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn1.Close()

	conn2, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn2.Close()

	assertCommandString(t, &conn1, "SET foo bar", "OK")
	assertCommandString(t, &conn2, "GET foo", "bar")
}

// Test many requests concurrently to try to break the server
func TestServer_ManyConnections(t *testing.T) {
	server := createServer()
	go server.Serve()
	time.Sleep(100 * time.Millisecond)

	numConnections := 100
	conns := make([]net.Conn, numConnections)
	for i := 0; i < numConnections; i++ {
		conn, err := net.Dial("tcp", "localhost:6379")
		if err != nil {
			t.Fatalf("Failed to connect to server: %v", err)
		}
		conns[i] = conn
		defer conn.Close()
	}

	for i := 0; i < numConnections; i++ {
		assertCommandString(t, &conns[i], "SET foo bar", "OK")
	}
}

func TestServer_StringOperations(t *testing.T) {
	server := createServer()
	go server.Serve()
	//defer server.Close() (check this later)
	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// APPEND existing key
	assertCommandString(t, &conn, "SET foo bar", "OK")
	assertCommandString(t, &conn, "APPEND foo baz", "6")
	assertCommandString(t, &conn, "GET foo", "barbaz")
	// APPEND non-existing key
	assertCommandString(t, &conn, "APPEND bar baz", "3")
	assertCommandString(t, &conn, "GET bar", "baz")
	// STRLEN existing key
	assertCommandString(t, &conn, "SET foo bar", "OK")
	assertCommandString(t, &conn, "STRLEN foo", "3")
	// STRLEN non-existing key
	assertCommandString(t, &conn, "STRLEN nonexistent", "0")
}
