package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/SXsid/kitsuDB/internal/config"
)

var conCount int

func handleConn(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("error while closing the connection ")
		}
		conCount -= 1
		log.Println("user", conn.RemoteAddr(), "disconnected. concurrent:", conCount)
	}()
	for {
		buff := make([]byte, 1024)

		n, err := conn.Read(buff)
		if err != nil {
			if err != io.EOF {
				log.Printf("error while reading the command \n error:%v", err)
			}
			break
		}
		res := buff[:n]
		_, err = conn.Write(res)
		fmt.Println(string(res))
		if err != nil {
			log.Printf("error while responding to client %s \b error:%v", conn.RemoteAddr(), err)
			break
		}
	}
}

func Run() {
	// tcp server
	listner, err := net.Listen("tcp", config.Cnfg.Host+":"+strconv.Itoa(config.Cnfg.Port))
	if err != nil {
		log.Fatalf("error while starting the server %v", err)
	}
	log.Println("kitsu is ready to eat data on ", config.Cnfg.Host, ":", config.Cnfg.Port)
	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Println("error while accepting new connection \n error", err.Error())
		}
		conCount += 1
		log.Println("clinet connectec with address", conn.RemoteAddr(), "concurent client", conCount)

		handleConn(conn)

	}
}
