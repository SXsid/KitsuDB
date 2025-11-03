package handler_test

import (
	"io"
	"testing"

	"github.com/SXsid/kitsuDB/internal/Handler"
)

func TestRespondWithError(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		err  error
		conn io.ReadWriter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.RespondWithError(tt.err, tt.conn)
		})
	}
}
