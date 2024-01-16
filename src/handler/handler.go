package handler

import "github.com/brkss/redis/src/resp"

var Handlers = map[string]func([]resp.Value) resp.Value{
	"PING": ping,
}

func ping(args []resp.Value) resp.Value {
	return resp.Value{Typ: "string", Str: "PONG"}
}
