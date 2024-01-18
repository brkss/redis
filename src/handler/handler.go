package handler

import (
	"sync"

	"github.com/brkss/redis/src/resp"
)

// SETs is a global variable to hold database data
var SETs = map[string]string{}

// SETsMutex is a global variable that hold mutexed to handle editing SETs's data
var SETsMutex = sync.RWMutex{}

// Handlers is global variable that hold function for each command !
var Handlers = map[string]func([]resp.Value) resp.Value{
	"PING":    ping,
	"GET":     get,
	"SET":     set,
	"HGET":    hget,
	"HSET":    hset,
	"HGETALL": hgetall,
}

func ping(args []resp.Value) resp.Value {
	return resp.Value{Typ: "string", Str: "PONG"}
}

func set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of argements for 'set' command"}
	}

	key := args[0].Blk
	val := args[1].Blk

	SETsMutex.Lock()
	SETs[key] = val
	SETsMutex.Unlock()

	return resp.Value{Typ: "string", Str: "OK"}
}

func get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Blk

	SETsMutex.RLock()
	val, ok := SETs[key]
	SETsMutex.RUnlock()

	if !ok {
		return resp.Value{Typ: "null"}
	}

	return resp.Value{Typ: "bulk", Blk: val}
}
