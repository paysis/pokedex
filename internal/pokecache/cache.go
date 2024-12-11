package pokecache

import (
	"bytes"
	"sync"
	"time"
)

type Cache struct {
	entries  map[string]cacheEntry
	lock     *sync.Mutex
	ticker   *time.Ticker
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       *bytes.Buffer
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entries:  map[string]cacheEntry{},
		lock:     &sync.Mutex{},
		ticker:   time.NewTicker(interval),
		interval: interval,
	}
	go cache.reapLoop()
	return cache
}

func (c *Cache) Add(key string, val *bytes.Buffer) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(), // as the app is local, lets use local version of time
		val:       val,
	}
}

func (c *Cache) Get(key string) (*bytes.Buffer, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	entry, ok := c.entries[key]
	return entry.val, ok
}

func (c *Cache) reapLoop() {
	defer c.ticker.Stop()
	for currentTime := range c.ticker.C {
		c.lock.Lock()
		for k, v := range c.entries {
			if currentTime.Sub(v.createdAt) >= c.interval {
				delete(c.entries, k)
			}
		}
		c.lock.Unlock()
	}
}
