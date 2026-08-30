package cache

import "time"

type Item struct {
	value      string
	expiration int64
}

func (item *Item) isExpired() bool {
	if item.expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.expiration
}
