package main

import (
	"fmt"
	"sync"
	"time"
)

type Item struct {
	value      string
	expiration int64
}

type Cache struct {
	mu    sync.RWMutex
	store map[string]Item
}

func NewCache(cleanupInterval time.Duration) *Cache {
	c := &Cache{
		store: make(map[string]Item),
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

	c.store[key] = Item{
		value:      value,
		expiration: exp,
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, exists := c.store[key]
	if !exists {
		return "", false
	}

	if val.isExpired() {
		delete(c.store, key)
		return "", false
	}

	return val.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.store, key)
}

func (c *Cache) startActiveCleaner(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.mu.Lock()
		now := time.Now().UnixNano()
		for k, item := range c.store {
			if item.expiration > 0 && now > item.expiration {
				delete(c.store, k)
			}
		}
		c.mu.Unlock()
	}
}

func main() {
	cache := NewCache(500 * time.Millisecond)

	cache.Set("session_token", "ABC123XYZ", 2*time.Second)

	if val, ok := cache.Get("session_token"); ok {
		fmt.Printf("[0s] Found key: %s\n", val)
	}

	time.Sleep(1 * time.Second)
	if val, ok := cache.Get("session_token"); ok {
		fmt.Printf("[1s] Found key: %s\n", val)
	}

	time.Sleep(1500 * time.Millisecond)
	if _, ok := cache.Get("session_token"); !ok {
		fmt.Println("[2.5s] Key successfully expired and cleaned up!")
	}
}
