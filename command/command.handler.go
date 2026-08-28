package command

import (
	"fmt"
	"gocache/pkg/cache"
)

type Engine struct {
	cache *cache.Cache
}

func NewEngine(c *cache.Cache) *Engine {
	return &Engine{
		cache: c,
	}
}

func (e *Engine) Executed(cmd Command) (string, error) {
	switch cmd.Name {
	case "SET":
		return e.handleSet(cmd)
	case "GET":
		return e.handleGet(cmd)
	case "DEL":
		return e.handleDelete(cmd)
	case "EXISTS":
		return e.handleExists(cmd)
	// case "EXPIRE":
	// 	return e.handleExpire(cmd)
	// case "TTL":
	// 	return e.handleTTL(cmd)
	default:
		return "", fmt.Errorf("unknown commad: %s", cmd.Name)
	}
}

func (e *Engine) handleSet(cmd Command) (string, error) {
	if len(cmd.Args) != 2 {
		return "", fmt.Errorf("SET requires key and value")
	}

	key := cmd.Args[0]
	value := cmd.Args[1]

	e.cache.Set(key, value, 0)

	return "OK", nil
}

func (e *Engine) handleGet(cmd Command) (string, error) {
	if len(cmd.Args) != 1 {
		return "", fmt.Errorf("GET requires key")
	}

	value, exists := e.cache.Get(cmd.Args[0])
	if !exists {
		return "(nil)", nil
	}

	return value, nil
}

func (e *Engine) handleDelete(cmd Command) (string, error) {
	if len(cmd.Args) == 0 {
		return "", fmt.Errorf("DEL requires at least one key")
	}
	count := e.cache.Delete(cmd.Args...)

	return fmt.Sprintf("%d", count), nil
}

func (e *Engine) handleExists(cmd Command) (string, error) {
	if len(cmd.Args) == 0 {
		return "", fmt.Errorf("EXISTS requires at least one key")
	}

	result := e.cache.Exists(cmd.Args...)
	count := 0

	for _, exists := range result {
		if exists {
			count++
		}
	}

	return fmt.Sprintf("%d", count), nil
}
