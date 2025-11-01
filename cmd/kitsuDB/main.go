package main

import (
	"flag"
	"fmt"

	resp "github.com/SXsid/kitsuDB/RESP"
	"github.com/SXsid/kitsuDB/internal/config"
	"github.com/SXsid/kitsuDB/internal/server"
)

func setUpFlags() {
	flag.StringVar(&config.Cnfg.Host, "host", "0.0.0.0", "address of your fox")
	flag.IntVar(&config.Cnfg.Port, "port", 8080, "room of your fox")
	flag.Parse()
}

func main() {
	setUpFlags()
	fmt.Println("🦊 Kitsu is waking up!")
	fmt.Println(resp.Decoder([]byte("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n")))
	server.Run()
}
