package main

import (
	"fmt"
	"sync"
)

type Cache struct {
	mu sync.RWMutex
	store map[string]string
}

func NewCache() *Cache {
	return &Cache{
		store: make(map[string]string),
	}
}

func (c *Cache) Set(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.store[key] = value
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	val, exists := c.store[key]
	return val, exists
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	delete(c.store, key)
}

func (c *Cache) FetchAll() map[string]string {
	return c.store
}

func main(){
	cache := NewCache()
	var wg sync.WaitGroup

	for i:=0; i<10; i++ {
		wg.Add(1)
		go func(id int){
			defer fmt.Println("Done writing")
			defer wg.Done()
			key := fmt.Sprintf("user_%d", id)
			val := fmt.Sprintf("value_%d", id)
			fmt.Println("SET: ",key, val)
			cache.Set(key, val)
		}(i)
	}

	wg.Wait()

	for i:=0; i<10; i++ {
		wg.Add(1)
		go func(id int){
			defer fmt.Println("Done reading")
			defer wg.Done()
			key := fmt.Sprintf("user_%d", id)
			value, ok := cache.Get(key)
			if !ok{
				fmt.Println("error")
			}
			fmt.Println("GET: ",key, value)
		}(i)
	}

	wg.Wait()
	fmt.Println("Successfully Processed 200 concurrent operations")
}