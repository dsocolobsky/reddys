package main

import (
	"github.com/dsocolobsky/reddys/internal"
	"testing"
)

func TestPing(t *testing.T) {
	handler := NewHandler(internal.NewMapDatabase())
	resp := handler.HandleCommand([]string{"ping"})
	if resp != "+PONG\r\n" {
		t.Errorf("Expected '+PONG\\r\\n', got '%s'", resp)
	}
}

func TestPingWithArguments(t *testing.T) {
	handler := NewHandler(internal.NewMapDatabase())
	resp := handler.HandleCommand([]string{"ping", "hello"})
	if resp != "$5\r\nhello\r\n" {
		t.Errorf("Expected '$5\\r\\nhello\\r\\n', got '%s'", resp)
	}
}

func TestSetAndGet(t *testing.T) {
	handler := NewHandler(internal.NewMapDatabase())
	resp := handler.HandleCommand([]string{"get", "key"})
	if resp != "_\r\n" {
		t.Errorf("Expected '_\\r\\n', got '%s'", resp)
	}
	resp = handler.HandleCommand([]string{"set", "key", "value"})
	if resp != "+OK\r\n" {
		t.Errorf("Expected '+OK\\r\\n', got '%s'", resp)
	}
	resp = handler.HandleCommand([]string{"get", "key"})
	if resp != "$5\r\nvalue\r\n" {
		t.Errorf("Expected '$5\\r\\nvalue\\r\\n', got '%s'", resp)
	}
}

func TestMgetAndMSet(t *testing.T) {
	handler := NewHandler(internal.NewMapDatabase())
	resp := handler.HandleCommand([]string{"mset", "key1", "val1", "key2", "val2"})
	if resp != "+OK\r\n" {
		t.Errorf("Expected '+OK\\r\\n', got '%s'", resp)
	}
	resp = handler.HandleCommand([]string{"mget", "key1", "key2", "key3"})
	if resp != "*3\r\n$4\r\nval1\r\n$4\r\nval2\r\n_\r\n" {
		t.Errorf("Expected '*3\\r\\n$4\\r\\nval1\\r\\n$4\\r\\nval2\\r\\n_\\r\\n', got '%s'", resp)
	}
}

func TestIncrAndDecr(t *testing.T) {
	handler := NewHandler(internal.NewMapDatabase())
	resp := handler.HandleCommand([]string{"incr", "key"})
	if resp != ":1\r\n" {
		t.Errorf("Expected ':1\\r\\n', got '%s'", resp)
	}
	resp = handler.HandleCommand([]string{"incr", "key"})
	if resp != ":2\r\n" {
		t.Errorf("Expected ':2\\r\\n', got '%s'", resp)
	}
	resp = handler.HandleCommand([]string{"decr", "key"})
	if resp != ":1\r\n" {
		t.Errorf("Expected ':1\\r\\n', got '%s'", resp)
	}
	resp = handler.HandleCommand([]string{"decr", "otherkey"})
	if resp != ":-1\r\n" {
		t.Errorf("Expected ':-1\\r\\n', got '%s'", resp)
	}
	resp = handler.HandleCommand([]string{"decr", "otherkey"})
	if resp != ":-2\r\n" {
		t.Errorf("Expected ':-2\\r\\n', got '%s'", resp)
	}
}

func TestIncrByAndDecrBy(t *testing.T) {
	handler := NewHandler(internal.NewMapDatabase())
	resp := handler.HandleCommand([]string{"incrby", "key", "5"})
	if resp != ":5\r\n" {
		t.Errorf("Expected ':5\\r\\n', got '%s'", resp)
	}
	resp = handler.HandleCommand([]string{"incrby", "key", "5"})
	if resp != ":10\r\n" {
		t.Errorf("Expected ':10\\r\\n', got '%s'", resp)
	}
	resp = handler.HandleCommand([]string{"decrby", "key", "3"})
	if resp != ":7\r\n" {
		t.Errorf("Expected ':7\\r\\n', got '%s'", resp)
	}
	resp = handler.HandleCommand([]string{"decrby", "key", "3"})
	if resp != ":4\r\n" {
		t.Errorf("Expected ':4\\r\\n', got '%s'", resp)
	}
}
