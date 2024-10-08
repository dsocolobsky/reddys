package server

import (
	"fmt"
	"github.com/dsocolobsky/reddys/pkg/database"
	"github.com/dsocolobsky/reddys/pkg/resp"
	"net"
	"strconv"
	"strings"
)

// Commands that should write to disk
var writeCommands = map[string]bool{
	"SET":    true,
	"MSET":   true,
	"HSET":   true,
	"INCR":   true,
	"DECR":   true,
	"INCRBY": true,
	"DECRBY": true,
}

type Handler struct {
	database database.Database
	aof      *database.AOF
}

func NewHandler(database database.Database, aof *database.AOF) *Handler {
	return &Handler{
		database: database,
		aof:      aof,
	}
}

func (h *Handler) ping(commands []string) string {
	if len(commands) == 1 {
		return resp.CraftSimpleString("PONG")
	} else if len(commands) == 2 {
		return resp.CraftBulkString(commands[1])
	}
	return resp.CraftSimpleError("wrong number of arguments for 'ping' command")
}

func (h *Handler) get(commands []string) string {
	if len(commands) != 2 {
		return resp.CraftSimpleError("wrong number of arguments for 'get' command")
	}
	key := commands[1]
	h.database.Lock()
	value := h.database.Get(key)
	h.database.Unlock()
	return resp.CraftAppropiateString(value)
}

func (h *Handler) set(commands []string) string {
	if len(commands) != 3 {
		return resp.CraftSimpleError("wrong number of arguments for 'set' command")
	}
	key := commands[1]
	value := commands[2]
	h.database.Lock()
	h.database.Set(key, value)
	h.database.Unlock()
	return resp.CraftSimpleString("OK")
}

func (h *Handler) mget(commands []string) string {
	if len(commands) < 2 {
		return resp.CraftSimpleError("wrong number of arguments for 'mget' command")
	}
	var values []string
	h.database.Lock()
	defer h.database.Unlock()
	for i := 1; i < len(commands); i++ {
		key := commands[i]
		value := h.database.Get(key)
		var st string
		if value == "" {
			st = resp.CraftNullString()
		} else {
			st = resp.CraftBulkString(value)
		}
		values = append(values, st)
	}
	return resp.CraftArray(values)
}

func (h *Handler) mset(commands []string) string {
	if len(commands) < 3 || len(commands)%2 != 1 {
		return resp.CraftSimpleError("wrong number of arguments for 'mset' command")
	}
	h.database.Lock()
	defer h.database.Unlock()
	for i := 1; i < len(commands); i += 2 {
		key := commands[i]
		value := commands[i+1]
		h.database.Set(key, value)
	}
	return resp.CraftSimpleString("OK")
}

func (h *Handler) hget(commands []string) string {
	if len(commands) != 3 {
		return resp.CraftSimpleError("wrong number of arguments for 'hget' command")
	}
	key := commands[1]
	field := commands[2]
	h.database.Lock()
	value := h.database.HGet(key, field)
	h.database.Unlock()
	return resp.CraftAppropiateString(value)
}

func (h *Handler) hset(commands []string) string {
	if len(commands) != 4 {
		return resp.CraftSimpleError("wrong number of arguments for 'hset' command")
	}
	key := commands[1]
	field := commands[2]
	value := commands[3]
	h.database.Lock()
	h.database.HSet(key, field, value)
	h.database.Unlock()
	return resp.CraftSimpleString("OK")
}

func (h *Handler) incr(commands []string) string {
	if len(commands) != 2 {
		return resp.CraftSimpleError("wrong number of arguments for 'incr' command")
	}
	return h._intModifyBy(commands[1], 1)
}

func (h *Handler) decr(commands []string) string {
	if len(commands) != 2 {
		return resp.CraftSimpleError("wrong number of arguments for 'decr' command")
	}
	return h._intModifyBy(commands[1], -1)
}

func (h *Handler) incrBy(commands []string) string {
	if len(commands) != 3 {
		return resp.CraftSimpleError("wrong number of arguments for 'incrby' command")
	}
	amount, err := strconv.Atoi(commands[2])
	if err != nil {
		return resp.CraftSimpleError("value is not an integer or out of range")
	}
	return h._intModifyBy(commands[1], amount)
}

func (h *Handler) decrBy(commands []string) string {
	if len(commands) != 3 {
		return resp.CraftSimpleError("wrong number of arguments for 'decrby' command")
	}
	amount, err := strconv.Atoi(commands[2])
	if err != nil {
		return resp.CraftSimpleError("value is not an integer or out of range")
	}
	return h._intModifyBy(commands[1], -amount)
}

func (h *Handler) _intModifyBy(key string, amount int) string {
	var intVal int
	h.database.Lock()
	defer h.database.Unlock()
	value := h.database.Get(key)
	if value == "" {
		intVal = 0
		h.database.Set(key, strconv.Itoa(intVal))
	} else {
		var err error
		intVal, err = strconv.Atoi(value)
		if err != nil {
			return resp.CraftSimpleError("value is not an integer or out of range")
		}
	}
	intVal += amount
	h.database.Set(key, strconv.Itoa(intVal))
	return resp.CraftInteger(intVal)
}

func (h *Handler) HandleCommand(commands []string) string {
	command := strings.ToUpper(strings.TrimSpace(commands[0]))

	switch command {
	case "PING":
		return h.ping(commands)
	case "SET":
		return h.set(commands)
	case "GET":
		return h.get(commands)
	case "MGET":
		return h.mget(commands)
	case "MSET":
		return h.mset(commands)
	case "HGET":
		return h.hget(commands)
	case "HSET":
		return h.hset(commands)
	case "INCR":
		return h.incr(commands)
	case "DECR":
		return h.decr(commands)
	case "INCRBY":
		return h.incrBy(commands)
	case "DECRBY":
		return h.decrBy(commands)
	}
	return resp.CraftSimpleError("ERR unknown command")
}

func Serve() {
	handler := NewHandler(database.NewMapDatabase(), database.NewAOF("database.aof"))
	defer handler.aof.Close()

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
		go handleConn(conn, handler)
	}
}

func toCommandArray(message string) []string {
	arr := resp.ReadArray(message)
	if len(arr) == 0 {
		fmt.Println("Empty array of commands!")
	}
	return arr
}

func handleConn(conn net.Conn, handler *Handler) {
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
		arr := toCommandArray(msg)
		if len(arr) > 0 && writeCommands[strings.ToUpper(arr[0])] {
			handler.aof.Write(msg)
		}
		res := handler.HandleCommand(arr)
		_, err = conn.Write([]byte(res))
		if err != nil {
			panic(err)
		}
	}
}
