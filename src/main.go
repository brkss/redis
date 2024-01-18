package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/brkss/redis/src/aof"
	"github.com/brkss/redis/src/handler"
	"github.com/brkss/redis/src/resp"
)

// Port is a global const representing the server port
const (
	Port       = ":6379"
	BufferSize = 1024
	AOFFile    = "./database.aof"
)

func main() {

	fmt.Println("Start listening for connection on 127.0.0.1" + Port)

	// create new server !
	l, err := net.Listen("tcp", Port)
	if err != nil {
		log.Fatal("Can't run redos server :", err)
	}

	// listen for connection
	conn, err := l.Accept()
	if err != nil {
		log.Fatal("Can't Accept connection : ", err)
	}

	// init aof
	aof, err := aof.NewAOF(AOFFile)
	if err != nil {
		log.Fatal("Something went wrong opening aof : ", err)
	}

	//read from old file
	err = aof.Read(func(val resp.Value) {
		command := strings.ToUpper(val.Arr[0].Blk)
		args := val.Arr[1:]

		handler, ok := handler.Handlers[command]
		if !ok {
			fmt.Println("invalid command : ", command)
			return
		}
		handler(args)
	})

	fmt.Println("read error : ", err)

	defer aof.Close()
	defer conn.Close()

	for {
		reader := resp.NewReader(conn)
		val, err := reader.Read()
		if err != nil {
			log.Fatal("Something went wrong reading client query : ", err)
		}

		if val.Typ != "array" {
			fmt.Println("request expected array !")
			continue
		}

		if len(val.Arr) == 0 {
			fmt.Println("expected array length to be greater than 0")
			continue
		}

		writer := resp.NewWriter(conn)

		command := strings.ToUpper(val.Arr[0].Blk)
		args := val.Arr[1:]

		handler, ok := handler.Handlers[command]
		if !ok {
			fmt.Println("invalid command : ", command)
			writer.Write(resp.Value{Typ: "string", Str: ""})
			continue
		}

		if command == "SET" || command == "HSET" {
			err := aof.Write(val)
			if err != nil {
				log.Println("something went wrong sync data : ", err)
			}
		}

		results := handler(args)
		writer.Write(results)
	}

}
