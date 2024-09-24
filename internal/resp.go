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
func ReadBulkString(message string) (string, int) {
	if message[0] != '$' {
		panic("Invalid bulk string")
	}
	message = message[1:]
	lineBreakIdx := strings.Index(message, "\r\n")
	if lineBreakIdx == -1 {
		panic("Invalid bulk string")
	}
	lengthStr := message[:lineBreakIdx]
	message = message[lineBreakIdx+2:]
	length := 0
	fmt.Sscanf(lengthStr, "%d", &length)
	if len(message) < length {
		panic("Invalid bulk string")
	}
	message = message[:length]
	// Add 4 for both pairs of \r\n, one for $ and len(lengthStr) for the number in the length
	return message[:length], length + 4 + 1 + len(lengthStr)
}

func ReadArray(message string) []string {
	fmt.Println("Reading array: ", message)
	if message[0] != '*' {
		panic("Invalid array")
	}
	message = message[1:]
	firstLineBreakIdx := strings.Index(message, "\r\n")
	if firstLineBreakIdx == -1 {
		panic("Invalid array")
	}
	lengthStr := message[:firstLineBreakIdx]
	message = message[firstLineBreakIdx+2:]
	length := 0
	fmt.Sscanf(lengthStr, "%d", &length)
	fmt.Println("Array length: ", length)
	array := make([]string, length)
	arrayIdx := 0
	for len(message) > 0 {
		fmt.Println("msg: ", message)
		msg, read := ReadRESP(message)
		array[arrayIdx] = msg
		message = message[read:]
		arrayIdx++
	}
	return array
}

func ReadSimpleString(message string) (string, int) {
	return readSimple(message, "+")
}

func ReadSimpleError(message string) (string, int) {
	return readSimple(message, "-")
}

func readSimple(message string, ch string) (string, int) {
	if message[0] != ch[0] {
		panic("Invalid simple string")
	}
	message = message[1:]
	// Here we stop at the first \r\n and ignore everything that follows, might be wrong.
	splitted := strings.SplitN(message, "\r\n", 2)
	return splitted[0], len(splitted[0]) + 3 // Add 3 to account for the + and \r\n
}

func ReadRESP(message string) (string, int) {
	message = strings.TrimSpace(message)
	switch message[0] {
	case '+':
		return ReadSimpleString(message)
	case '-':
		return ReadSimpleError(message)
	case '$':
		return ReadBulkString(message)
	case '*':
		// Not yet implemented
		//return ReadArray(message)
	default:
		panic("Invalid RESP message")
	}
	return "", -1
}
