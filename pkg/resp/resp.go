package resp

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

// CraftAppropiateString crafts a RESP bulk string if message is not empty, otherwise it crafts a null string
func CraftAppropiateString(message string) string {
	if message == "" {
		return CraftNullString()
	}
	return CraftBulkString(message)
}

func CraftNullString() string {
	return "_\r\n"
}

// CraftBoolean crafts a RESP boolean of the type "#t\r\n" or "#f\r\n"
func CraftBoolean(val bool) string {
	if val {
		return "#t\r\n"
	}
	return "#f\r\n"
}

// CraftInteger crafts a RESP integer of the type ":(+|-)?val\r\n"
func CraftInteger(val int) string {
	return fmt.Sprintf(":%d\r\n", val)
}

func CraftArray(array []string) string {
	nElems := len(array)
	return "*" + fmt.Sprintf("%d\r\n", nElems) + strings.Join(array, "")
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

func ReadArray(message string) ([]string, int) {
	totalRead := 0
	if message[0] != '*' {
		panic("Invalid array")
	}
	message = message[1:]
	totalRead += 1
	firstLineBreakIdx := strings.Index(message, "\r\n")
	if firstLineBreakIdx == -1 {
		panic("Invalid array")
	}
	numElemsStr := message[:firstLineBreakIdx]
	message = message[firstLineBreakIdx+2:]
	totalRead += firstLineBreakIdx + 2
	numElems := 0
	fmt.Sscanf(numElemsStr, "%d", &numElems)
	array := make([]string, numElems)
	arrayIdx := 0
	for arrayIdx < numElems {
		msg, read := ReadRESP(message)
		array[arrayIdx] = msg
		totalRead += read
		message = message[read:]
		arrayIdx++
	}
	return array, totalRead
}

func ReadManyArrays(message string) [][]string {
	arrays := make([][]string, 0)
	for len(message) > 0 {
		arr, read := ReadArray(message)
		arrays = append(arrays, arr)
		message = message[read:]
	}
	return arrays
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

func ReadBoolean(message string) (string, int) {
	// TODO We should return a boolean here, but for now let's do string.
	fmt.Println("Boolean: ", message)
	if len(message) < 4 || message[0] != '#' {
		panic("Invalid boolean")
	}
	val := message[1]
	if val == 't' {
		return "true", 4
	} else if val == 'f' {
		return "false", 4
	}
	panic("Invalid boolean, no t/f")
}

func ReadInteger(message string) (string, int) {
	// TODO we should return integer here instead of string
	if message[0] != ':' {
		panic("Invalid integer")
	}
	message = message[1:]
	var sign byte
	// Handle optional sign
	if message[0] == '-' || message[0] == '+' {
		sign = message[0]
		message = message[1:]
	}
	splitted := strings.SplitN(message, "\r\n", 2)
	intStr := splitted[0]
	if sign == '-' {
		return "-" + intStr, len(intStr) + 4 // Add 2 for the \r\n and 2 for the : and -
	} else if sign == '+' {
		return intStr, len(intStr) + 4 // Add 2 for the \r\n and 2 for the : and +
	}
	return intStr, len(intStr) + 3 // Add 2 for the \r\n and 1 for the :
}

func ReadRESP(message string) (string, int) {
	switch message[0] {
	case '+':
		return ReadSimpleString(message)
	case '-':
		return ReadSimpleError(message)
	case '$':
		return ReadBulkString(message)
	case '#':
		return ReadBoolean(message)
	case ':':
		return ReadInteger(message)
	case '*':
		// Not yet implemented
		//return ReadArray(message)
	default:
		panic("Invalid RESP message")
	}
	return "", -1
}
