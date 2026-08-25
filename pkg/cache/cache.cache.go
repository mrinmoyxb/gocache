package cache

import (
	"sync"
	"time"
)

type Cache struct {
	mu    sync.RWMutex
	store map[string]Item
}

func NewCache(cleanupInterval time.Duration) *Cache {
	c := &Cache{
		store: make(map[string]Item),
	}

	if cleanupInterval > 0 {
		go c.startSweeper(cleanupInterval)
	}

	return c
}
