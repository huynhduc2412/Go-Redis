package core

import (
	"Go-Redis/internal/constant"
	"errors"
	"fmt"
	"strconv"
	"syscall"
	"time"
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

func cmdSet(args []string) []byte {
	if len(args) < 2 || len(args) == 3 || len(args) > 4 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SET' commad") , false)
	}
	var key , value string
	var ttlMs int64 = -1
	key , value = args[0] , args[1]
	if len(args) > 2 {
		ttlSec , err := strconv.ParseInt(args[3] , 10 , 64)
		if err != nil {
			return Encode(errors.New("(error) ERR value is not an integer or out of range") , false)
		}
		ttlMs = ttlSec * 1000
	}
	dicStore.Set(key , dicStore.NewObj(key , value , ttlMs))
	return constant.RespOk
}

func cmdGet(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'GET' command") , false)
	}
	key := args[0]
	obj := dicStore.Get(key)
	if obj == nil {
		return constant.RespNil
	}

	if dicStore.HasExpired(key) {
		return constant.RespNil
	}
	return Encode(obj.Value , false)
}

func cmdTTL(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'TTL' command") , false)
	}
	key := args[0]
	obj := dicStore.Get(key)
	if obj == nil {
		return constant.TtlKeyNotExist
	}

	exp , isExpirySet := dicStore.GetExpiry(key)
	if !isExpirySet {
		return constant.TtlKeyExistNoExpire
	}
	
	remainMs := int64(exp) - int64(time.Now().UnixMilli())
	if remainMs < 0 {
		return constant.TtlKeyNotExist
	}
	return Encode(int64(remainMs / 1000) , false)
}

// ExecuteAndResponse given a Command, executes it and responses
func ExecuteAndResponse(cmd *Command , connFd int) error {
	var res []byte
	switch cmd.Cmd {
	case "PING":
		res = cmdPING(cmd.Args)
	case "SET": 
		res = cmdSet(cmd.Args)
	case "GET":
		res = cmdGet(cmd.Args)
	case "TTL":
		res = cmdTTL(cmd.Args)
	default:
		res = []byte(fmt.Sprintf("-CMD NOT FOUND\r\n"))
	}
	_, err := syscall.Write(connFd , res)
	return err
}