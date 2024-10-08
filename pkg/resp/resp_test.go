package resp

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

func TestCraftBooleanTrue(t *testing.T) {
	expected := "#t\r\n"
	actual := CraftBoolean(true)
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestCraftBooleanFalse(t *testing.T) {
	expected := "#f\r\n"
	actual := CraftBoolean(false)
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

func TestCraftInteger(t *testing.T) {
	expected := ":-123\r\n"
	actual := CraftInteger(-123)
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

func TestReadBoolean(t *testing.T) {
	expected := "true"
	raw := "#t\r\n"
	actual, read := ReadBoolean(raw)
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
	if read != len(raw) {
		t.Errorf("Expected to read %d but got %d", len(raw), read)
	}
}

func TestReadInteger(t *testing.T) {
	expected := "123"
	raw := ":123\r\n"
	actual, read := ReadInteger(raw)
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
	if read != len(raw) {
		t.Errorf("Expected to read %d but got %d", len(raw), read)
	}
}

func TestReadArrayOfSimpleStrings(t *testing.T) {
	expected := []string{"one", "two", "three"}
	str := "*3\r\n+one\r\n+two\r\n+three\r\n"
	actual, read := ReadArray(str)
	fmt.Println(expected)
	fmt.Println(actual)
	for i, v := range expected {
		if v != actual[i] {
			t.Errorf("Expected %s but got %s", v, actual[i])
		}
	}
	if read != len(str) {
		t.Errorf("Expected to read %d but got %d", len(str), read)
	}
}

func TestReadArrayOfBulkStrings(t *testing.T) {
	expected := []string{"one", "two", "three"}
	str := "*3\r\n$3\r\none\r\n$3\r\ntwo\r\n$5\r\nthree\r\n"
	actual, read := ReadArray(str)
	fmt.Println(expected)
	fmt.Println(actual)
	for i, v := range expected {
		if v != actual[i] {
			t.Errorf("Expected %s but got %s", v, actual[i])
		}
	}
	if read != len(str) {
		t.Errorf("Expected to read %d but got %d", len(str), read)
	}
}

func TestReadArrayOfManyBulkStrings(t *testing.T) {
	expected := []string{"one", "two", "three"}
	str1 := "*3\r\n$3\r\none\r\n$3\r\ntwo\r\n$5\r\nthree\r\n"
	str2 := "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	actual, read := ReadArray(str1 + str2)
	fmt.Println(expected)
	fmt.Println(actual)
	for i, v := range expected {
		if v != actual[i] {
			t.Errorf("Expected %s but got %s", v, actual[i])
		}
	}
	if read != len(str1) {
		t.Errorf("Expected to read %d but got %d", len(str1), read)
	}
}

func TestReadManyArrays(t *testing.T) {
	expected := [][]string{
		{"one", "two", "three"},
		{"four", "five", "six"},
	}
	actual := ReadManyArrays("*3\r\n+one\r\n+two\r\n+three\r\n*3\r\n+four\r\n+five\r\n+six\r\n")
	for i, arr := range expected {
		for j, v := range arr {
			if v != actual[i][j] {
				t.Errorf("Expected %s but got %s", v, actual[i][j])
			}
		}
	}
}

func TestReadArrayOfBooleans(t *testing.T) {
	expected := []string{"true", "false", "true"}
	str := "*3\r\n#t\r\n#f\r\n#t\r\n"
	actual, read := ReadArray(str)
	for i, v := range expected {
		if v != actual[i] {
			t.Errorf("Expected %s but got %s", v, actual[i])
		}
	}
	if read != len(str) {
		t.Errorf("Expected to read %d but got %d", len(str), read)
	}
}

func TestReadArrayOfIntegers(t *testing.T) {
	expected := []string{"10", "-9", "7"}
	str := "*3\r\n:10\r\n:-9\r\n:+7\r\n"
	actual, read := ReadArray(str)
	for i, v := range expected {
		if v != actual[i] {
			t.Errorf("Expected %s but got %s", v, actual[i])
		}
	}
	if read != len(str) {
		t.Errorf("Expected to read %d but got %d", len(str), read)
	}
}
