package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

const (
	PORT        = ":6379"
	BUFFER_SIZE = 1024
	DELIMITER
)

func main() {

	fmt.Println("Start listening for connection on 127.0.0.1" + PORT)

	// create new server !
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal("Can't run redos server :", err)
	}

	// listen for connection
	conn, err := l.Accept()
	if err != nil {
		log.Fatal("Can't Accept connection : ", err)
	}

	defer conn.Close()

	for {
		buff := make([]byte, BUFFER_SIZE)
		_, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal("Something went wrong reading client query : ", err)
		}

		conn.Write([]byte("+OK\r\n"))
	}

}
