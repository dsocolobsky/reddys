package internal

import (
	"fmt"
	"strings"
)

// CraftSimpleString crafts a RESP simple string of the type "+message\r\n"
func CraftSimpleString(message string) string {
	return "+" + message + "\r\n"
}

// CraftSimpleError crafts a RESP simple error of the type "-message\r\n"
func CraftSimpleError(message string) string {
	return "-" + message + "\r\n"
}

// CraftBulkString crafts a RESP bulk string of the type "$length\r\nmessage\r\n"
func CraftBulkString(message string) string {
	length := fmt.Sprintf("%d", len(message))
	return "$" + length + "\r\n" + message + "\r\n"
}

// ReadBulkString reads a RESP bulk string of the type "$length\r\nmessage\r\n" into a string "message"
func ReadBulkString(message string) string {
	if message[0] != '$' {
		panic("Invalid bulk string")
	}
	message = message[1:]
	lineBreakIdx := strings.Index(message, "\r\n")
	if lineBreakIdx == -1 {
		panic("Invalid bulk string")
	}
	length := message[:lineBreakIdx]
	message = message[lineBreakIdx+2:]
	lengthInt := 0
	fmt.Sscanf(length, "%d", &lengthInt)
	if len(message) < lengthInt {
		panic("Invalid bulk string")
	}
	return message[:lengthInt]
}
