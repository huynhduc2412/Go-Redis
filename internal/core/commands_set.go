package core

import (
	datastructure "Go-Redis/internal/data_structure"
	"errors"
)

func cmdSADD(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SADD' command") , false)
	}
	//key each set
	key := args[0]
	set , exist := setStore[key]
	if !exist {
		set = datastructure.NewSimpleSet()
		setStore[key] = set
	}
	count := set.Add(args[1:]...)
	return Encode(count , false)
}

func cmdSREM(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SREM' command"), false)
	}
	key := args[0]
	set , exist := setStore[key]
	if !exist {
		set = datastructure.NewSimpleSet()
		setStore[key] = set
	}
	count := set.Rem(args[1:]...)
	return Encode(count , false)
}

func cmdSMEMBERS(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SMEMBERS' command"), false)
	} 
	key := args[0]
	set , exist := setStore[key]
	if !exist {
		return Encode(make([]string , 0) , false)
	}
	return Encode(set.Members() , false)
}

func cmdSISMEMBER(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SISMEMBER' command"), false)
	}
	key := args[0]
	set , exist := setStore[key]
	if !exist {
		return Encode(0 , false)
	}
	return Encode(set.IsMember(args[1]) , false)
}