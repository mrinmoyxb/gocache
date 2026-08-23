package test

import (
	"fmt"
	"sync"
)

type Items struct {
	value string
	expiration int64
}

type Cache struct {
	mw sync.RWMutex
	store map[string]Item
}

func NewCache(cleanupInterval time.Duration) *Cache {
	c := Cache {
		store : make(map[string]Item),
	}

	if cleanupInterval > 0 {
		go c.startActiveCleaner(cleanupInterval)
	}

	return c
}

func (item Item) isExpired() bool {
	if item.expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.expiration
}

func (c *Cache) Set(key string, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	c.store[key] = Item {
		value : value,
		expiration: exp,
	}
}

func (c *Cache) st
