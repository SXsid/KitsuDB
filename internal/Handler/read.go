package handler

import (
	"io"
	"strings"

	decode "github.com/SXsid/kitsuDB/internal/Decode"
	"github.com/SXsid/kitsuDB/internal/config"
)

func ReadCommand(conn io.ReadWriter) (*config.Input, error) {
	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		return nil, err
	}
	data, err := decode.DecodeArrayString(buff[:n])
	if err != nil {
		return nil, err
	}
	return &config.Input{
		Command: strings.ToUpper(data[0]),
		Args:    data[1:],
	}, nil
}
