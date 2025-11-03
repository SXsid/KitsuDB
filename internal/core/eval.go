package core

import (
	"errors"

	"github.com/SXsid/kitsuDB/internal/config"
)

func Eval(cmd *config.Input) ([]byte, error) {
	switch cmd.Command {
	case "PING":
		return PING(cmd.Args)
	}
	return nil, errors.New("not a valid commnad")
}

func PING(args []string) ([]byte, error) {
	if len(args) >= 2 {
		return nil, errors.New("wrong number of arguments for 'ping' command")
	}
	if len(args) == 0 {
		return Encode("PONG", true), nil
	} else {
		return Encode(args[0], false), nil
	}
}
