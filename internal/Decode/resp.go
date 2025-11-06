package decode

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

func Decoder(data []byte) (any, error) {
	if len(data) == 0 {
		return nil, errors.New("no data found")
	}
	value, _, err := Parser(data)
	return value, err
}

// int is the parsing count
func Parser(data []byte) (any, int, error) {
	switch data[0] {
	case '*':
		return ReadArray(data)
	case ':':
		return ReadNumber(data)
	case '$':
		return ReadBulkString(data)
	case '+':
		return ReadSimpleString(data)
	case '-':
		return ReadError(data)
	}
	return nil, 0, nil
}

func DecodeArrayString(data []byte) ([]string, error) {
	value, err := Decoder(data)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, fmt.Errorf("'%s' is not a valid command", string(data))
	}
	// type cast in array of any
	result_array := value.([]any)
	Array_strings := make([]string, len(result_array))
	for i := range result_array {
		Array_strings[i] = result_array[i].(string)
	}
	return Array_strings, nil
}

func ReadArray(data []byte) ([]any, int, error) {
	// delta is the pos where the content length ended form 1 to first\r\n
	count, delta, err := ReadNumber(data)
	pos := int64(delta)
	if err != nil {
		return nil, 0, err
	}
	values := make([]any, count)
	for i := range values {
		item, delta, err := Parser(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		values[i] = item
		pos += int64(delta)

	}

	return values, int(pos), nil
}

func ReadBulkString(data []byte) (string, int, error) {
	number, delta, err := ReadNumber(data)
	pos := int64(delta)
	if err != nil {
		return "", 0, err
	}
	return string(data[pos:(pos + number)]), (delta + int(number) + 2), nil
}

func ReadError(data []byte) (string, int, error) {
	pos, err := findPos(data)
	if err != nil {
		return "", 0, err
	}
	return string(data[1:pos]), pos + 2, nil
}

func findPos(data []byte) (int, error) {
	pos := bytes.IndexByte(data, '\r')
	if pos == -1 {
		return 0, errors.New("no CRLF present")
	}
	return pos, nil
}

func ReadNumber(data []byte) (int64, int, error) {
	pos, err := findPos(data)
	if err != nil {
		return 0, 0, err
	}
	num, err := strconv.ParseInt(string(data[1:pos]), 10, 64)
	if err != nil {
		return 0, 0, errors.New("invalid number")
	}
	return num, pos + 2, nil
}

func ReadSimpleString(data []byte) (string, int, error) {
	return ReadError(data)
}
