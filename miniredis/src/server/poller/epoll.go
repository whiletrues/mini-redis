package poller

import (
	"syscall"
)

// File descriptor for the epoll instance
// File descriptor for the eventfd used for waking up the epoll instance
type Poller struct {
	fd      int
	eventFd int
}

// Create creates a new Poller instance.
// It returns a pointer to the Poller and an error, if any.
// The Poller is responsible for managing epoll events.
func Create() (*Poller, error) {
	fd, err := syscall.EpollCreate1(0)
	if err != nil {
		return nil, err
	}

	eventFdPtr, _, err := syscall.Syscall(syscall.SYS_EVENTFD2, 0, 0, 0)

	if err != nil {
		syscall.Close(fd)
		return nil, err
	}

	eventFd := int(eventFdPtr)

	err = syscall.EpollCtl(fd, syscall.EPOLL_CTL_ADD, eventFd, &syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(eventFd),
	})

	if err != nil {
		syscall.Close(fd)
		syscall.Close(eventFd)
		return nil, err
	}

	return &Poller{
		fd:      fd,
		eventFd: eventFd,
	}, nil

}

// close closes the poller by closing the eventFd and fd.
// It returns an error if there was an error while closing the file descriptors.
func (poller *Poller) Close() (err error) {
	err = syscall.Close(poller.eventFd)

	if err != nil {
		return err
	}
	err = syscall.Close(poller.fd)

	if err != nil {
		return err
	}

	return
}

// Register adds a file descriptor to the epoll instance for monitoring events.
// It returns an error if the registration fails.
// The fd parameter specifies the file descriptor to be registered.
func (poller *Poller) AddRead(fd int) (err error) {
	return syscall.EpollCtl(poller.fd, syscall.EPOLL_CTL_ADD, fd, &syscall.EpollEvent{
		Events: syscall.EPOLLIN | syscall.EPOLLPRI,
		Fd:     int32(fd),
	})
}

// AddWrite adds a file descriptor to the epoll instance for write events.
// It registers the file descriptor with the EPOLLIN and EPOLLPRI events.
// The fd parameter specifies the file descriptor to be added.
// It returns an error if the operation fails.
func (poller *Poller) AddWrite(fd int) (err error) {
	return syscall.EpollCtl(poller.fd, syscall.EPOLL_CTL_ADD, fd, &syscall.EpollEvent{
		Events: syscall.EPOLLIN | syscall.EPOLLPRI,
		Fd:     int32(fd),
	})
}

// Delete removes the file descriptor from the epoll instance.
// It returns an error if the operation fails.
func (poller *Poller) Delete(fd int) (err error) {
	return syscall.EpollCtl(poller.fd, syscall.EPOLL_CTL_DEL, fd, nil)
}

// Update updates the events for a file descriptor in the epoll instance.
// It modifies the events associated with the given file descriptor.
// The fd parameter specifies the file descriptor to update.
// The events parameter specifies the new events to associate with the file descriptor.
// It returns an error if the update operation fails.
func (poller *Poller) Update(fd int, events uint32) (err error) {
	return syscall.EpollCtl(poller.fd, syscall.EPOLL_CTL_MOD, fd, &syscall.EpollEvent{
		Events: events,
		Fd:     int32(fd),
	})
}

func (poller *Poller) Poll(handler func(fd int, event uint32)) (int, error) {

	events := make([]syscall.EpollEvent, 1024)

	for {

		nevents, err := syscall.EpollWait(poller.fd, events, 1000)

		if err != nil && err != syscall.EINTR {
			continue
		}

		for ev := 0; ev < nevents; ev++ {
			if int(events[ev].Fd) != poller.eventFd {
				syscall.Read(poller.eventFd, []byte{0, 0, 0, 0, 0, 0, 0, 0})

				handler(int(events[ev].Fd), events[ev].Events)

			} else {
				println("wake")
			}
		}
	}
}
