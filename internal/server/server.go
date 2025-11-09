package server

import (
	"fmt"
	"log"
	"net"

	handler "github.com/SXsid/kitsuDB/internal/Handler"
	"github.com/SXsid/kitsuDB/internal/config"
	"golang.org/x/sys/unix"
)

func Run() error {
	con_client := 0
	max_client := 20 * 1000
	serverFD, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM|unix.SOCK_NONBLOCK, 0)
	if err != nil {
		return err
	}
	defer unix.Close(serverFD)
	if err := unix.SetsockoptInt(serverFD, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); err != nil {
		return err
	}

	ipv4 := net.ParseIP(config.Cnfg.Host).To4()
	add := unix.SockaddrInet4{
		Port: config.Cnfg.Port,
		Addr: [4]byte{ipv4[0], ipv4[1], ipv4[2], ipv4[3]},
	}
	if err := unix.Bind(serverFD, &add); err != nil {
		return err
	}
	if err := unix.Listen(serverFD, max_client); err != nil {
		return err
	}
	events := make([]unix.EpollEvent, max_client)
	epoll_fd, err := unix.EpollCreate1(0)
	if err != nil {
		return err
	}
	defer unix.Close(epoll_fd)
	serverEvent := &unix.EpollEvent{
		Fd:     int32(serverFD),
		Events: unix.EPOLLIN,
	}
	if err := unix.EpollCtl(epoll_fd, unix.EPOLL_CTL_ADD, serverFD, serverEvent); err != nil {
		return err
	}
	fmt.Println("Fox is up and hungry on ", config.Cnfg.Host, ":", config.Cnfg.Port)
	for {
		n, err := unix.EpollWait(epoll_fd, events, -1)
		if err != nil {
			continue
		}
		for i := 0; i < n; i++ {
			eventFD := events[i].Fd
			if eventFD == int32(serverFD) {
				for {
					fd, _, err := unix.Accept(serverFD)
					if err != nil {
						if err == unix.EAGAIN || err == unix.EWOULDBLOCK {
							break // no more clients waiting
						}
						log.Println("error while connecting to a client")
						continue
					}
					if err := unix.SetNonblock(fd, true); err != nil {
						return err
					}
					con_client++
					clientEvent := unix.EpollEvent{
						Fd:     int32(fd),
						Events: unix.EPOLLIN,
					}
					if err := unix.EpollCtl(epoll_fd, unix.EPOLL_CTL_ADD, fd, &clientEvent); err != nil {
						log.Println("err", err)
						continue
					}
					log.Println("a clinet connected", con_client)
				}
			} else {
				conn := config.Conn{
					Fd: int(eventFD),
				}
				cmd, err := handler.ReadCommand(conn)
				if err != nil {
					handler.RespondWithError(err, conn)
					unix.Close(int(eventFD))

					con_client -= 1
					log.Println("connection losed", con_client)

					continue
				}
				handler.Respond(cmd, conn)
			}
		}
	}
}
