package command

import (
	"gocache/pkg/cache"
	"gocache/protocol"
)

type Engine struct {
	cache *cache.Cache
}

func NewEngine(c *cache.Cache) *Engine {
	return &Engine{
		cache: c,
	}
}

func (e *Engine) Executed(cmd Command) (protocol.Response, error) {
	switch cmd.Name {
	case "SET":
		return e.handleSet(cmd)
	case "GET":
		return e.handleGet(cmd)
	case "DEL":
		return e.handleDelete(cmd)
	case "EXISTS":
		return e.handleExists(cmd)
	default:
		return protocol.Error("unknown command"), nil
	}
}

func (e *Engine) handleSet(cmd Command) (protocol.Response, error) {
	if len(cmd.Args) != 2 {
		return protocol.Error("wrong number of arguments for 'SET'"), nil
	}

	key := cmd.Args[0]
	value := cmd.Args[1]

	e.cache.Set(key, value, 0)

	return protocol.SimpleString("OK"), nil
}

func (e *Engine) handleGet(cmd Command) (protocol.Response, error) {
	if len(cmd.Args) != 1 {
		return protocol.Error("wrong number of arguments for 'GET'"), nil
	}

	value, exists := e.cache.Get(cmd.Args[0])
	if !exists {
		return protocol.Null(), nil
	}

	return protocol.BulkString(value), nil
}

func (e *Engine) handleDelete(cmd Command) (protocol.Response, error) {
	if len(cmd.Args) == 0 {
		return protocol.Error("wrong number of arguments for 'GET'"), nil
	}
	count := e.cache.Delete(cmd.Args...)

	return protocol.Integer(int64(count)), nil
}

func (e *Engine) handleExists(cmd Command) (protocol.Response, error) {
	if len(cmd.Args) == 0 {
		return protocol.Error("wrong number of arguments for 'EXISTS'"), nil
	}

	result := e.cache.Exists(cmd.Args...)
	count := 0

	for _, exists := range result {
		if exists {
			count++
		}
	}

	return protocol.Integer(int64(count)), nil
}
