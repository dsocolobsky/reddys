package main

import "testing"

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
	actual := ReadBulkString("$7\r\nmessage\r\n")
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestReadBulkStringMultipleLines(t *testing.T) {
	expected := "some\nother\r\nmessage"
	actual := ReadBulkString("$19\r\nsome\nother\r\nmessage\r\n")
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}
