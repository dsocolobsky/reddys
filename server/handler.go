package main

import (
	"github.com/dsocolobsky/reddys/internal"
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
	value := h.database.Get(key)
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
	h.database.Set(key, value)
	return internal.CraftSimpleString("OK")
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
	}
	return internal.CraftSimpleError("ERR unknown command")
}
