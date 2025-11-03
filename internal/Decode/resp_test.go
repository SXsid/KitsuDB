package decode

import (
	"reflect"
	"testing"
)

func TestReadNumber(t *testing.T) {
	tests := []struct {
		input    []byte
		expected int64
		consumed int
		hasErr   bool
	}{
		{[]byte(":123\r\n"), 123, 6, false},
		{[]byte(":-99\r\n"), -99, 6, false},
		{[]byte(":0\r\n"), 0, 4, false},
		{[]byte(":abc\r\n"), 0, 0, true},
	}

	for _, tt := range tests {
		n, c, err := ReadNumber(tt.input)
		if (err != nil) != tt.hasErr {
			t.Errorf("expected err=%v got=%v", tt.hasErr, err)
		}
		if n != tt.expected {
			t.Errorf("expected number %d got %d", tt.expected, n)
		}
		if c != tt.consumed {
			t.Errorf("expected consumed %d got %d", tt.consumed, c)
		}
	}
}

func TestReadSimpleString(t *testing.T) {
	input := []byte("+OK\r\n")
	value, consumed, err := ReadSimpleString(input)
	if err != nil {
		t.Fatal(err)
	}
	if value != "OK" {
		t.Errorf("expected OK got %s", value)
	}
	if consumed != 5 {
		t.Errorf("expected consumed 5 got %d", consumed)
	}
}

func TestReadError(t *testing.T) {
	input := []byte("-ERR something\r\n")
	value, consumed, err := ReadError(input)
	if err != nil {
		t.Fatal(err)
	}
	if value != "ERR something" {
		t.Errorf("expected ERR something got %s", value)
	}
	if consumed != len(input) {
		t.Errorf("unexpected consumed size")
	}
}

func TestReadBulkString(t *testing.T) {
	input := []byte("$5\r\nhello\r\n")
	value, consumed, err := ReadBulkString(input)
	if err != nil {
		t.Fatal(err)
	}
	if value != "hello" {
		t.Errorf("expected hello got %s", value)
	}
	if consumed != len(input) {
		t.Errorf("unexpected consumed size")
	}
}

func TestReadArray(t *testing.T) {
	input := []byte("*2\r\n:1\r\n:2\r\n")
	value, consumed, err := ReadArray(input)
	if err != nil {
		t.Fatal(err)
	}

	expected := []any{int64(1), int64(2)}
	if !reflect.DeepEqual(value, expected) {
		t.Errorf("expected %v got %v", expected, value)
	}
	if consumed != len(input) {
		t.Errorf("unexpected consumed size")
	}
}

func TestDecoderFull(t *testing.T) {
	input := []byte("*3\r\n:1\r\n+OK\r\n$5\r\nhello\r\n")

	value, err := Decoder(input)
	if err != nil {
		t.Fatal(err)
	}

	expected := []any{
		int64(1),
		"OK",
		"hello",
	}

	if !reflect.DeepEqual(value, expected) {
		t.Errorf("expected %+v got %+v", expected, value)
	}
}
