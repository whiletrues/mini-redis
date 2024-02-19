package epoll

import (
	"golang.org/x/sys/unix"
)

const (
	EPOLLIN      = unix.EPOLLIN
	EPOLLOUT     = unix.EPOLLOUT
	EPOLLRDHUP   = unix.EPOLLRDHUP
	EPOLLPRI     = unix.EPOLLPRI
	EPOLLERR     = unix.EPOLLERR
	EPOLLHUP     = unix.EPOLLHUP
	EPOLLET      = unix.EPOLLET
	EPOLLONESHOT = unix.EPOLLONESHOT

	// _EPOLLCLOSED is a special EpollEvent value the receipt of which means
	// that the epoll instance is closed.
	EPOLLCLOSED = 0x20
)

type EpollEvent uint32

// File descriptor for the epoll instance
// File descriptor for the eventfd used for waking up the epoll instance
type EPoll struct {
	fd        int
	eventFd   int
	callbacks map[int]func(EpollEvent)
}

// Create creates a new Poller instance.
// It returns a pointer to the Poller and an error, if any.
// The Poller is responsible for managing epoll events.
func Create() (*EPoll, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}

	eventFdPtr, _, errno := unix.Syscall(unix.SYS_EVENTFD2, 0, 0, 0)
	if errno != 0 {
		return nil, errno
	}
	eventFd := int(eventFdPtr)

	err = unix.EpollCtl(fd, unix.EPOLL_CTL_ADD, eventFd, &unix.EpollEvent{
		Events: unix.EPOLLIN,
		Fd:     int32(eventFd),
	})

	if err != nil {
		unix.Close(fd)
		unix.Close(eventFd)
		return nil, err
	}

	return &EPoll{
		fd:        fd,
		eventFd:   eventFd,
		callbacks: make(map[int]func(EpollEvent)),
	}, nil

}

// close closes the poller by closing the eventFd and fd.
// It returns an error if there was an error while closing the file descriptors.
func (poller *EPoll) Close() (err error) {
	err = unix.Close(poller.eventFd)

	if err != nil {
		return err
	}
	err = unix.Close(poller.fd)

	if err != nil {
		return err
	}

	for _, cb := range poller.callbacks {
		cb(EPOLLCLOSED)
	}
	return
}

// Register adds a file descriptor to the epoll instance for monitoring events.
// It returns an error if the registration fails.
// The fd parameter specifies the file descriptor to be registered.
func (poller *EPoll) AddRead(fd int, events EpollEvent, cb func(EpollEvent)) (err error) {

	poller.callbacks[fd] = cb

	return unix.EpollCtl(poller.fd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{
		Events: uint32(events),
		Fd:     int32(fd),
	})
}

// Delete removes the file descriptor from the epoll instance.
// It returns an error if the operation fails.
func (poller *EPoll) Delete(fd int) (err error) {
	return unix.EpollCtl(poller.fd, unix.EPOLL_CTL_DEL, fd, nil)
}

// Update updates the events for a file descriptor in the epoll instance.
// It modifies the events associated with the given file descriptor.
// The fd parameter specifies the file descriptor to update.
// The events parameter specifies the new events to associate with the file descriptor.
// It returns an error if the update operation fails.
func (poller *EPoll) Update(fd int, events uint32) (err error) {
	return unix.EpollCtl(poller.fd, unix.EPOLL_CTL_MOD, fd, &unix.EpollEvent{
		Events: events,
		Fd:     int32(fd),
	})
}

func (poller *EPoll) Poll(onError func(error)) {

	defer func() {
		if err := unix.Close(poller.fd); err != nil {
			onError(err)
		}
	}()

	events := make([]unix.EpollEvent, 1024)
	callbacks := make([]func(EpollEvent), 0, 32768)

	for {

		n, err := unix.EpollWait(poller.fd, events, -1)

		if err != nil && err != unix.EINTR {
			onError(err)
		}

		callbacks = callbacks[:n]

		for i := 0; i < n; i++ {
			fd := int(events[i].Fd)
			if fd == poller.eventFd {
				return
			}
			callbacks[i] = poller.callbacks[fd]
		}

		for i := 0; i < n; i++ {
			if cb := callbacks[i]; cb != nil {
				cb(EpollEvent(events[i].Events))
				callbacks[i] = nil
			}
		}

		if n == len(events) && n*2 <= 32768 {
			events = make([]unix.EpollEvent, n*2)
			callbacks = make([]func(EpollEvent), 0, n*2)
		}
	}
}
