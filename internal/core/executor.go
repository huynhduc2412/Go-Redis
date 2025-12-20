package core

import (
	"errors"
	"fmt"
	"syscall"
)

func cmdPING(args []string) []byte {
	var res []byte
	if len(args) > 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'ping' command") , false)
	}
	if len(args) == 0 {
		return Encode("PONG" , true)
	}else {
		res = Encode(args[0] , false)
	}
	return res
}

func ExecuteAndResponse(cmd *Command , connFd int) error {
	var res []byte
	switch cmd.Cmd {
	case "PING":
		res = cmdPING(cmd.Args)
	default:
		res = []byte(fmt.Sprintf("-CMD NOT FOUND\r\n"))
	}
	_, err := syscall.Write(connFd , res)
	return err
}