package main

import (
	"github.com/dsocolobsky/reddys/internal"
	"strconv"
	"strings"
)

type Handler struct {
	database internal.Database
}

func NewHandler(database internal.Database) *Handler {
	return &Handler{
		database: database,
	}
}

func (h *Handler) ping(commands []string) string {
	if len(commands) == 1 {
		return internal.CraftSimpleString("PONG")
	} else if len(commands) == 2 {
		return internal.CraftBulkString(commands[1])
	}
	return internal.CraftSimpleError("wrong number of arguments for 'ping' command")
}

func (h *Handler) get(commands []string) string {
	if len(commands) != 2 {
		return internal.CraftSimpleError("wrong number of arguments for 'get' command")
	}
	key := commands[1]
	h.database.Lock()
	value := h.database.Get(key)
	h.database.Unlock()
	if value == "" {
		return internal.CraftNullString()
	} else {
		return internal.CraftBulkString(value)
	}
}

func (h *Handler) set(commands []string) string {
	if len(commands) != 3 {
		return internal.CraftSimpleError("wrong number of arguments for 'set' command")
	}
	key := commands[1]
	value := commands[2]
	h.database.Lock()
	h.database.Set(key, value)
	h.database.Unlock()
	return internal.CraftSimpleString("OK")
}

func (h *Handler) mget(commands []string) string {
	if len(commands) < 2 {
		return internal.CraftSimpleError("wrong number of arguments for 'mget' command")
	}
	var values []string
	h.database.Lock()
	defer h.database.Unlock()
	for i := 1; i < len(commands); i++ {
		key := commands[i]
		value := h.database.Get(key)
		var st string
		if value == "" {
			st = internal.CraftNullString()
		} else {
			st = internal.CraftBulkString(value)
		}
		values = append(values, st)
	}
	return internal.CraftArray(values)
}

func (h *Handler) mset(commands []string) string {
	if len(commands) < 3 || len(commands)%2 != 1 {
		return internal.CraftSimpleError("wrong number of arguments for 'mset' command")
	}
	h.database.Lock()
	defer h.database.Unlock()
	for i := 1; i < len(commands); i += 2 {
		key := commands[i]
		value := commands[i+1]
		h.database.Set(key, value)
	}
	return internal.CraftSimpleString("OK")
}

func (h *Handler) incr(commands []string) string {
	if len(commands) != 2 {
		return internal.CraftSimpleError("wrong number of arguments for 'incr' command")
	}
	return h._intModifyBy(commands[1], 1)
}

func (h *Handler) decr(commands []string) string {
	if len(commands) != 2 {
		return internal.CraftSimpleError("wrong number of arguments for 'decr' command")
	}
	return h._intModifyBy(commands[1], -1)
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
			return internal.CraftSimpleError("value is not an integer or out of range")
		}
	}
	intVal += amount
	h.database.Set(key, strconv.Itoa(intVal))
	return internal.CraftInteger(intVal)
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
	case "INCR":
		return h.incr(commands)
	case "DECR":
		return h.decr(commands)
	}
	return internal.CraftSimpleError("ERR unknown command")
}
