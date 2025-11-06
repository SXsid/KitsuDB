package config

import "golang.org/x/sys/unix"

type Conn struct {
	Fd int
}

func (conn Conn) Write(b []byte) (int, error) {
	return unix.Write(conn.Fd, b)
}

func (conn Conn) Read(b []byte) (int, error) {
	return unix.Read(conn.Fd, b)
}
