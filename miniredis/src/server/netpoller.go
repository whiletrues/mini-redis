package server

import (
	epoll "miniredis/src/server/poller"
	"os"
)

type eventLoop struct {
	*epoll.EPoll
}

type Event uint16

const (
	// EventHup is indicates that some side of i/o operations (receive, send or
	// both) is closed.
	// Usually (depending on operating system and its version) the EventReadHup
	// or EventWriteHup are also set int Event value.
	EventHup Event = 0x10

	EventReadHup  = 0x20
	EventWriteHup = 0x40

	EventErr = 0x80

	// EventPollerClosed is a special Event value the receipt of which means that the
	// Poller instance is closed.
	EventPollerClosed = 0x8000
)

const (
	EventOneShot       Event = 0x4
	EventEdgeTriggered       = 0x8
)
const (
	EventRead  Event = 0x1
	EventWrite       = 0x2
)

func New() (*eventLoop, error) {

	poller, err := epoll.Create()

	if err != nil {
		return nil, err
	}

	return &eventLoop{poller}, nil

}

func (ev *eventLoop) Start(desc *Desc, cb func(Event)) error {

	err := ev.EPoll.AddRead(desc.fd(), toEpollEvent(desc.event),
		func(ep epoll.EpollEvent) {
			var event Event

			if ep&epoll.EPOLLHUP != 0 {
				event |= EventHup
			}
			if ep&epoll.EPOLLRDHUP != 0 {
				event |= EventReadHup
			}
			if ep&epoll.EPOLLIN != 0 {
				event |= EventRead
			}
			if ep&epoll.EPOLLOUT != 0 {
				event |= EventWrite
			}
			if ep&epoll.EPOLLERR != 0 {
				event |= EventErr
			}
			if ep&epoll.EPOLLCLOSED != 0 {
				event |= EventPollerClosed
			}

			cb(event)
		},
	)
	if err == nil {
		if err = setNonblock(desc.fd(), true); err != nil {
			return os.NewSyscallError("setnonblock", err)
		}
	}

	return err
}

func toEpollEvent(event Event) (ep epoll.EpollEvent) {
	if event&EventRead != 0 {
		ep |= epoll.EPOLLIN | epoll.EPOLLRDHUP
	}
	if event&EventWrite != 0 {
		ep |= epoll.EPOLLOUT
	}
	if event&EventOneShot != 0 {
		ep |= epoll.EPOLLONESHOT
	}
	if event&EventEdgeTriggered != 0 {
		ep |= epoll.EPOLLET
	}
	return ep
}
