package server

import (
	"log"
	"net"

	handler "github.com/SXsid/kitsuDB/internal/Handler"
	"github.com/SXsid/kitsuDB/internal/config"
	"golang.org/x/sys/unix"
)

var conClinet int

func Run_server() error {
	max_clinet := 20 * 1000
	// event buffer to keep track /notfiication
	events := make([]unix.EpollEvent, max_clinet)

	serverFD, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM|unix.SOCK_NONBLOCK, 0)
	if err != nil {
		return err
	}
	defer unix.Close(serverFD)

	ipv4 := net.ParseIP(config.Cnfg.Host).To4()
	add := &unix.SockaddrInet4{
		Port: config.Cnfg.Port,
		Addr: [4]byte{ipv4[0], ipv4[1], ipv4[2], ipv4[3]},
	}
	if err := unix.Bind(serverFD, add); err != nil {
		return err
	}
	if err := unix.Listen(serverFD, max_clinet); err != nil {
		return err
	}
	log.Println("kitsu is ready to eat data on ", config.Cnfg.Host, ":", config.Cnfg.Port)
	epollFD, err := unix.EpollCreate1(0)
	if err != nil {
		return err
	}
	defer unix.Close(epollFD)
	serverEvent := unix.EpollEvent{
		Fd:     int32(serverFD),
		Events: unix.EPOLLIN,
	}
	if err := unix.EpollCtl(epollFD, unix.EPOLL_CTL_ADD, serverFD, &serverEvent); err != nil {
		return err
	}

	for {
		availEventNumber, err := unix.EpollWait(epollFD, events, -1)
		if err != nil {
			// keep servering others
			continue
		}
		for i := 0; i < availEventNumber; i++ {
			if events[i].Fd == int32(serverFD) {
				for {
					fd, _, err := unix.Accept(serverFD)
					if err != nil {
						// accepted all pendening fd iin the bloackge queue
						if err == unix.EAGAIN {
							break
						}
						log.Println("err", err)
						continue
					}
					conClinet++
					unix.SetNonblock(fd, true)

					clientEvent := &unix.EpollEvent{
						Fd:     int32(fd),
						Events: unix.EPOLLIN,
					}
					if err := unix.EpollCtl(epollFD, unix.EPOLL_CTL_ADD, fd, clientEvent); err != nil {
						log.Printf("err:%s", err)

						return err
					}

				}
			} else {
				clientFd := events[i].Fd
				conn := config.Conn{
					Fd: int(clientFd),
				}
				cmd, err := handler.ReadCommand(conn)
				if err != nil {
					unix.Close(int(clientFd))
					conClinet--
					continue
				}
				handler.Respond(cmd, conn)

			}
		}
	}
}
