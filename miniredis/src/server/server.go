package server

import (
	"fmt"
	"net"
	"strconv"
	"syscall"
)

type server struct {
	listener net.Listener
}

type Socket struct {
	FileDescriptor int
}

func (socket Socket) Read(buffer []byte) (int, error) {
	if len(buffer) == 0 {
		return 0, nil
	}

	readCount, err :=
		syscall.Read(socket.FileDescriptor, buffer)
	if err != nil {
		readCount = 0
	}
	return readCount, err
}

func (socket Socket) Write(buffer []byte) (int, error) {
	Writecount, err := syscall.Write(socket.FileDescriptor, buffer)
	if err != nil {
		Writecount = 0
	}
	return Writecount, err
}

func (socket Socket) Close() error {
	return syscall.Close(socket.FileDescriptor)
}

func (socket *Socket) String() string {
	return strconv.Itoa(socket.FileDescriptor)
}

func listen(ip string, port int) (*Socket, error) {

	socket := &Socket{}

	socketFileDescriptor, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)

	socket.FileDescriptor = socketFileDescriptor
	if err != nil {
		return nil, fmt.Errorf("failed to create socket (%v)", err)
	}

	socketAddress := &syscall.SockaddrInet4{Port: port}
	copy(socketAddress.Addr[:], net.ParseIP(ip).To4())

	if err = syscall.Bind(socket.FileDescriptor, socketAddress); err != nil {
		return nil, fmt.Errorf("failed to bind socket (%v)", err)
	}

	if err = syscall.Listen(socket.FileDescriptor, syscall.SOMAXCONN); err != nil {
		return nil, fmt.Errorf("failed to listen on socket (%v)", err)
	}

	return socket, nil
}

func StartServer(ip string, port int) (err error) {
	socket, err := listen(ip, port)

	if err != nil {
		return fmt.Errorf("failed to start server (%v)", err)
	}

	epollDescriptor, err := syscall.EpollCreate1(0)

	if err != nil {
		return fmt.Errorf("epoll creation failed (%v)", err)
	}

	err = syscall.EpollCtl(epollDescriptor, syscall.EPOLL_CTL_ADD, socket.FileDescriptor, &syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(socket.FileDescriptor),
	})

	if err != nil {
		return fmt.Errorf("epoll ctl failed (%v)", err)
	}

	events := make([]syscall.EpollEvent, 10)

	for {
		nevents, err := syscall.EpollWait(epollDescriptor, events, -1)

		if err != nil {
			return fmt.Errorf("epoll wait failed (%v)", err)
		}
		for ev := 0; ev < nevents; ev++ {
			if int(events[ev].Fd) == int(socket.FileDescriptor) {

				connectionFileDescriptor, _, err := syscall.Accept(socket.FileDescriptor)

				if err != nil {
					return fmt.Errorf("socket accept failed (%v)", err)
				}

				println(connectionFileDescriptor)

			}
		}
	}
}

func handleConnection(connection *Socket) error {
	defer connection.Close()

	var buffer [1024]byte

	for {
		readCount, err := connection.Read(buffer[:])

		if err != nil {
			return fmt.Errorf("read failed (%v)", err)
		}

		if readCount == 0 {
			return nil
		}

		_, err = connection.Write(buffer[:readCount])

		if err != nil {
			return fmt.Errorf("write failed (%v)", err)
		}
	}
}
