package main

import (
	"io"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	log.Println(conn.RemoteAddr())
	//read data from clients
	for {
		cmd , err := readCommand(conn)
		if err != nil {
			conn.Close()
			log.Println("Client disconnected" , conn.RemoteAddr())
			if err == io.EOF {
				break
			}
			return
		}else{
			log.Println("Command:",cmd)
		}
		if err = respond(cmd , conn) ; err != nil {
			log.Println("error write:" , err)
		}
	}
}

func readCommand(c net.Conn) (string , error) {
	var buf []byte = make([]byte , 1000)
	n , err := c.Read(buf)
	if err != nil {
		return "" , err
	} 
	return string(buf[:n]) , err 
}

func respond(cmd string , c net.Conn) error {
	if _ , err := c.Write([]byte(cmd)) ; err != nil {
		return err
	}
	return nil
}

func main() {
	listener, err := net.Listen("tcp" , ":3000")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening at port 3000")
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