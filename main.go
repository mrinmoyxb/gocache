package main

import "fmt"

type Cache struct {
	store map[string]string
}

func NewCache() *Cache {
	return &Cache{
		store: make(map[string]string),
	}
}

func (c *Cache) Set(key string, value string) {
	c.store[key] = value
}

func (c *Cache) Get(key string) (string, bool) {
	val, exists := c.store[key]
	return val, exists
}

func (c *Cache) Delete(key string) {
	delete(c.store, key)
}

func (c *Cache) FetchAll() map[string]string {
	return c.store
}

func main() {
	
}
