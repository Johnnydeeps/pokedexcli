package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	jsonData  []byte
}

type Cache struct {
	cache map[string]cacheEntry
	mux   *sync.Mutex
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		cache: make(map[string]cacheEntry),
		mux:   &sync.Mutex{},
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		c.mux.Lock()
		cutoff := time.Now().UTC().Add(-interval)
		for key, value := range c.cache {
			if value.createdAt.Before(cutoff) {
				delete(c.cache, key)
			}
		}
		c.mux.Unlock()
	}
}

func (c *Cache) Add(key string, responseJson []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.cache[key] = cacheEntry{
		createdAt: time.Now().UTC(),
		jsonData:  responseJson,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	value, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	return value.jsonData, true
}
