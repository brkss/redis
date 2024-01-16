package main

import (
	"fmt"
	"log"
	"net"

	"github.com/brkss/redis/src/resp"
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

		reader := resp.NewReader(conn)
		val, err := reader.Read()
		if err != nil {
			log.Fatal("Something went wrong reading client query : ", err)
		}

		fmt.Println("val : ", val)
		writer := resp.NewWriter(conn)
		v := resp.Value{Typ: "string", Str: "woring.."}
		writer.Write(v)
		//conn.Write([]byte("+OK\r\n"))
	}

}
