package handler

import (
	"fmt"
	"io"
	"log"

	"github.com/SXsid/kitsuDB/internal/config"
	"github.com/SXsid/kitsuDB/internal/core"
)

func Respond(cmd *config.Input, conn io.ReadWriter) {
	res, err := core.Eval(cmd)
	if err != nil {
		RespondWithError(err, conn)
	}
	// respond to the command
	if _, err := conn.Write(res); err != nil {
		log.Println("Error while responding to client", err)
	}
}

func RespondWithError(err error, conn io.ReadWriter) {
	// encode it with RESP specs for error
	if _, err := fmt.Fprintf(conn, "- ERR %s\r\n", err); err != nil {
		log.Println("Error while responding to client", err)
	}
}
