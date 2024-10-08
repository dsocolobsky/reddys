package resp

import (
	"fmt"
	"testing"
)

func TestMarshalString(t *testing.T) {
	expected := "+message\r\n"
	actual := MarshalString("message")
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestMarshalError(t *testing.T) {
	expected := "-message\r\n"
	actual := MarshalError("message")
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestMarshalInteger2BooleanTrue(t *testing.T) {
	expected := "#t\r\n"
	actual := MarshalBoolean(true)
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestMarshalInteger2BooleanFalse(t *testing.T) {
	expected := "#f\r\n"
	actual := MarshalBoolean(false)
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestMarshalInteger2BulkString(t *testing.T) {
	expected := "$7\r\nmessage\r\n"
	actual := MarshalBulkString("message")
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestMarshalInteger2BulkStringNull(t *testing.T) {
	expected := "_\r\n"
	actual := MarshalBulkString("")
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestMarshalInteger2Integer(t *testing.T) {
	expected := ":-123\r\n"
	actual := MarshalInteger(-123)
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestMarshalArrayOfBulkStrings(t *testing.T) {
	expected := "*3\r\n$3\r\none\r\n$3\r\ntwo\r\n$5\r\nthree\r\n"
	actual := MarshalArrayOfBulkStrings([]string{"one", "two", "three"})
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestUnmarshalBulkString(t *testing.T) {
	expected := "message"
	raw := "$7\r\nmessage\r\n"
	actual, read := UnmarshalBulkString(raw)
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
	if read != len(raw) {
		t.Errorf("Expected to read %d but got %d", len(raw), read)
	}
}

func TestUnmarshalBulkStringMultipleLines(t *testing.T) {
	expected := "some\nother\r\nmessage"
	raw := "$19\r\nsome\nother\r\nmessage\r\n"
	actual, read := UnmarshalBulkString(raw)
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
	if read != len(raw) {
		t.Errorf("Expected to read %d but got %d", len(raw), read)
	}
}

func TestUnmarshalSimpleString(t *testing.T) {
	expected := "message"
	raw := "+message\r\n"
	actual, read := UnmarshalString(raw)
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
	if read != len(raw) {
		t.Errorf("Expected to read %d but got %d", len(raw), read)
	}
}

func TestUnmarshalSimpleError(t *testing.T) {
	expected := "message"
	raw := "-message\r\n"
	actual, read := UnmarshalError(raw)
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
	if read != len(raw) {
		t.Errorf("Expected to read %d but got %d", len(raw), read)
	}
}

func TestUnmarshalBoolean(t *testing.T) {
	expected := "true"
	raw := "#t\r\n"
	actual, read := UnmarshalBoolean(raw)
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
	if read != len(raw) {
		t.Errorf("Expected to read %d but got %d", len(raw), read)
	}
}

func TestUnmarshalInteger(t *testing.T) {
	expected := "123"
	raw := ":123\r\n"
	actual, read := UnmarshalInteger(raw)
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
	if read != len(raw) {
		t.Errorf("Expected to read %d but got %d", len(raw), read)
	}
}

func TestUnmarshalArrayOfSimpleStrings(t *testing.T) {
	expected := []string{"one", "two", "three"}
	str := "*3\r\n+one\r\n+two\r\n+three\r\n"
	actual, read := UnmarshalArray(str)
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

func TestUnmarshalArrayOfBulkStrings(t *testing.T) {
	expected := []string{"one", "two", "three"}
	str := "*3\r\n$3\r\none\r\n$3\r\ntwo\r\n$5\r\nthree\r\n"
	actual, read := UnmarshalArray(str)
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

func TestUnmarshalArrayOfManyBulkStrings(t *testing.T) {
	expected := []string{"one", "two", "three"}
	str1 := "*3\r\n$3\r\none\r\n$3\r\ntwo\r\n$5\r\nthree\r\n"
	str2 := "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	actual, read := UnmarshalArray(str1 + str2)
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

func TestUnmarshalManyArrays(t *testing.T) {
	expected := [][]string{
		{"one", "two", "three"},
		{"four", "five", "six"},
	}
	actual := UnmarshalManyArrays("*3\r\n+one\r\n+two\r\n+three\r\n*3\r\n+four\r\n+five\r\n+six\r\n")
	for i, arr := range expected {
		for j, v := range arr {
			if v != actual[i][j] {
				t.Errorf("Expected %s but got %s", v, actual[i][j])
			}
		}
	}
}

func TestUnmarshalArrayOfBooleans(t *testing.T) {
	expected := []string{"true", "false", "true"}
	str := "*3\r\n#t\r\n#f\r\n#t\r\n"
	actual, read := UnmarshalArray(str)
	for i, v := range expected {
		if v != actual[i] {
			t.Errorf("Expected %s but got %s", v, actual[i])
		}
	}
	if read != len(str) {
		t.Errorf("Expected to read %d but got %d", len(str), read)
	}
}

func TestUnmarshalArrayOfIntegers(t *testing.T) {
	expected := []string{"10", "-9", "7"}
	str := "*3\r\n:10\r\n:-9\r\n:+7\r\n"
	actual, read := UnmarshalArray(str)
	for i, v := range expected {
		if v != actual[i] {
			t.Errorf("Expected %s but got %s", v, actual[i])
		}
	}
	if read != len(str) {
		t.Errorf("Expected to read %d but got %d", len(str), read)
	}
}
