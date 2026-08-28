package cache

import (
	"strings"
	"time"
)

func (c *Cache) Set(key string, value string, ttl time.Duration) bool {
	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = Item{value: value, expiration: exp}
	return true
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	item, exists := c.store[key]
	if !exists {
		c.mu.RUnlock()
		return "", false
	}

	if !item.isExpired() {
		c.mu.RUnlock()
		return item.value, true
	}

	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists = c.store[key]
	if !exists {
		return "", false
	}

	if exists && item.isExpired() {
		delete(c.store, key)
		return "", false
	}

	return item.value, true
}

func (c *Cache) Delete(keys ...string) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := 0
	for _, key := range keys {
		if _, exists := c.store[key]; exists {
			delete(c.store, key)
			count++
		}
	}
	return count
}

func (c *Cache) Exists(keys ...string) map[string]bool {
	result := make(map[string]bool, len(keys))

	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, key := range keys {
		item, exists := c.store[key]
		if exists && !item.isExpired() {
			result[key] = true
		} else {
			result[key] = false
		}
	}
	return result
}

func (c *Cache) Keys(pattern string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var result []string
	for key, item := range c.store {
		if item.isExpired() {
			continue
		}

		if pattern == "*" {
			result = append(result, key)
		} else if strings.HasPrefix(pattern, "*") && strings.HasSuffix(pattern, "*") {
			sub := strings.Trim(pattern, "*")
			if strings.Contains(key, sub) {
				result = append(result, key)
			}
		} else if strings.HasSuffix(pattern, "*") {
			prefix := strings.TrimSuffix(pattern, "*")
			if strings.HasPrefix(key, prefix) {
				result = append(result, key)
			}
		} else if key == pattern {
			result = append(result, key)
		}
	}
	return result
}
