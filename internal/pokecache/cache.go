package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mu   *sync.Mutex
	data map[string]cacheEntry
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = cacheEntry{
		createdAt: time.Now().UTC(),
		val:       value,
	}

}

func (c *Cache) Get(key string) (cacheData []byte, exists bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	value, exists := c.data[key]
	if !exists {
		return nil, exists
	}
	return value.val, true
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		mu:   &sync.Mutex{},
		data: make(map[string]cacheEntry),
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.data {
			if time.Since(entry.createdAt) > interval {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}

}
