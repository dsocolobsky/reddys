package internal

import (
	"fmt"
	"testing"
)

func TestCraftSimpleString(t *testing.T) {
	expected := "+message\r\n"
	actual := CraftSimpleString("message")
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestCraftSimpleError(t *testing.T) {
	expected := "-message\r\n"
	actual := CraftSimpleError("message")
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestCraftBulkString(t *testing.T) {
	expected := "$7\r\nmessage\r\n"
	actual := CraftBulkString("message")
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestCraftBulkStringNull(t *testing.T) {
	expected := "$0\r\n\r\n"
	actual := CraftBulkString("")
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestReadBulkString(t *testing.T) {
	expected := "message"
	raw := "$7\r\nmessage\r\n"
	actual, read := ReadBulkString(raw)
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
	if read != len(raw) {
		t.Errorf("Expected to read %d but got %d", len(raw), read)
	}
}

func TestReadBulkStringMultipleLines(t *testing.T) {
	expected := "some\nother\r\nmessage"
	raw := "$19\r\nsome\nother\r\nmessage\r\n"
	actual, read := ReadBulkString(raw)
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
	if read != len(raw) {
		t.Errorf("Expected to read %d but got %d", len(raw), read)
	}
}

func TestReadSimpleString(t *testing.T) {
	expected := "message"
	raw := "+message\r\n"
	actual, read := ReadSimpleString(raw)
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
	if read != len(raw) {
		t.Errorf("Expected to read %d but got %d", len(raw), read)
	}
}

func TestReadSimpleError(t *testing.T) {
	expected := "message"
	raw := "-message\r\n"
	actual, read := ReadSimpleError(raw)
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
	if read != len(raw) {
		t.Errorf("Expected to read %d but got %d", len(raw), read)
	}
}

func TestReadArrayOfSimpleStrings(t *testing.T) {
	expected := []string{"one", "two", "three"}
	actual := ReadArray("*3\r\n+one\r\n+two\r\n+three\r\n")
	fmt.Println(expected)
	fmt.Println(actual)
	for i, v := range expected {
		if v != actual[i] {
			t.Errorf("Expected %s but got %s", v, actual[i])
		}
	}
}

func TestReadArrayOfBulkStrings(t *testing.T) {
	expected := []string{"one", "two", "three"}
	actual := ReadArray("*3\r\n$3\r\none\r\n$3\r\ntwo\r\n$5\r\nthree\r\n")
	fmt.Println(expected)
	fmt.Println(actual)
	for i, v := range expected {
		if v != actual[i] {
			t.Errorf("Expected %s but got %s", v, actual[i])
		}
	}
}
