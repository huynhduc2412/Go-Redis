//go:build windows
// +build windows

package io_multiplexing

/*
#cgo CFLAGS: -I${SRCDIR}
#cgo LDFLAGS: -lws2_32
#include "wepoll.h"
*/
import "C"

import "unsafe"

func (e Event) toNative() C.struct_epoll_event {
	var events C.uint32_t = C.EPOLLIN
	if e.Op == OpWrite {
		events = C.EPOLLOUT
	}

	var ev C.struct_epoll_event
	ev.events = events

	*(*C.uint64_t)(unsafe.Pointer(&ev.data[0])) = C.uint64_t(e.Fd)

	return ev
}

func createEvent(ep C.struct_epoll_event) Event {
	var op Operation = OpRead
	if ep.events == C.EPOLLOUT {
		op = OpWrite
	}

	fd := *(*C.uint64_t)(unsafe.Pointer(&ep.data[0]))

	return Event{
		Fd: int(fd),
		Op: op,
	}
}
