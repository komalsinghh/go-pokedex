package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	mu       sync.Mutex
	DataMap  map[string]cacheEntry
	Interval time.Duration
}

type cacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		DataMap:  make(map[string]cacheEntry),
		Interval: interval,
	}
	go cache.DeleteData(interval)
	return cache
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.DataMap[key] = cacheEntry{
		Val:       value,
		CreatedAt: time.Now(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	data, exists := c.DataMap[key]
	if !exists {
		return []byte{}, false
	}
	return data.Val, true
}

func (c *Cache) DeleteData(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.mu.Lock()
		for key, value := range c.DataMap {
			if time.Since(value.CreatedAt) > c.Interval {
				fmt.Printf("Deleting Data %s", c.DataMap[key])
				delete(c.DataMap, key)
			}
		}
		c.mu.Unlock()
	}
}
