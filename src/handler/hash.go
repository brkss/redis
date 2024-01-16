package handler

import (
	"sync"

	"github.com/brkss/redis/src/resp"
)

var HSETs = map[string]map[string]string{}
var HSETsMutext = sync.RWMutex{}

func hset(args []resp.Value) resp.Value {
	if len(args) != 3 {
		return resp.Value{Typ: "error", Str: "ERR invalid arguments for 'hset' command"}
	}

	hash := args[0].Blk
	key := args[1].Blk
	val := args[2].Blk

	HSETsMutext.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = val
	HSETsMutext.Unlock()

	return resp.Value{Typ: "string", Str: "OK"}
}

func hget(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR invalid arguments for 'hget' command"}
	}

	hash := args[0].Blk
	key := args[1].Blk

	HSETsMutext.Lock()
	val, ok := HSETs[hash][key]
	HSETsMutext.Unlock()
	if !ok {
		return resp.Value{Typ: "null"}
	}

	return resp.Value{Typ: "bulk", Blk: val}
}

func hgetall(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR invalid arguments for 'hget' command"}
	}

	hash := args[0].Blk

	HSETsMutext.Lock()
	hashVal, ok := HSETs[hash]
	HSETsMutext.Unlock()

	if !ok {
		return resp.Value{Typ: "null"}
	}

	response := resp.Value{Typ: "array", Arr: make([]resp.Value, 0)}
	for _, val := range hashVal {
		item := resp.Value{Typ: "bulk", Blk: val}
		response.Arr = append(response.Arr, item)
	}

	return response
}
