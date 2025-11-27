//go:build windows
// +build windows

package io_multiplexing

/*
#cgo CFLAGS: -I${SRCDIR}
#cgo LDFLAGS: -lws2_32
#include "wepoll.h"
*/
import "C"

import (
	"Go-Redis/internal/config"
    "fmt"
    "unsafe"
)

type Epoll struct {
    fd            C.HANDLE
    epollEvents   []C.struct_epoll_event
    genericEvents []Event
}

func CreateIOMultiplexer() (*Epoll, error) {
    ep := C.epoll_create1(0)
    if uintptr(unsafe.Pointer(ep)) == 0  {
        return nil, fmt.Errorf("wepoll: epoll_create1 failed")
    }

    return &Epoll{
        fd:            ep,
        epollEvents:   make([]C.struct_epoll_event, config.MaxConnection),
        genericEvents: make([]Event, config.MaxConnection),
    }, nil
}

func (ep *Epoll) Monitor(event Event) error {
    ne := event.toNative()

    ret := C.epoll_ctl(
        ep.fd,
        C.EPOLL_CTL_ADD,
        C.SOCKET(event.Fd),
        (*C.struct_epoll_event)(unsafe.Pointer(&ne)),
    )
    if ret < 0 {
        return fmt.Errorf("wepoll: epoll_ctl failed")
    }
    return nil
}

func (ep *Epoll) Wait() ([]Event, error) {
    n := C.epoll_wait(
        ep.fd,
        (*C.struct_epoll_event)(unsafe.Pointer(&ep.epollEvents[0])),
        C.int(len(ep.epollEvents)),
        C.int(-1),
    )

    if n < 0 {
        return nil, fmt.Errorf("wepoll: epoll_wait failed")
    }

    for i := 0; i < int(n); i++ {
        ep.genericEvents[i] = createEvent(ep.epollEvents[i])
    }

    return ep.genericEvents[:n], nil
}

func (ep *Epoll) Close() error {
    C.epoll_close(ep.fd)
    return nil
}
