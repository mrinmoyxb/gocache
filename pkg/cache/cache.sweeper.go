package cache

import "time"

func (c *Cache) startSweeper(interval time.Duration){
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.mu.Lock()
		now := time.Now().UnixNano()
		for k, v := range c.store {
			if v.expiration > 0 && now > v.expiration {
				delete(c.store, k)
			}
		}
		c.mu.Unlock()
	}
}