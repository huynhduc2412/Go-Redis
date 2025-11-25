package main

import (
	"log"
	"net"
	"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close() 
	log.Println(conn.RemoteAddr())
	//read data from clients
	var buf []byte = make([]byte , 1000)
	for {
		_ , err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		//process
		time.Sleep(time.Second * 10)
		
		//rep 
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello, world\r\n"))
	}
}


func main() {
	listener, err := net.Listen("tcp" , ":3000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		// socket == dedicated communication channel
		conn , err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		//create new goroutine to handle the connection
		//thread per connection
		go handleConnection(conn)
	}
}