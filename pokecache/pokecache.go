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
    interval time.Duration
    data     map[string]cacheEntry
    mutex    sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
    cache := &Cache{
        interval: interval,
        data:     make(map[string]cacheEntry),
    }
    go cache.reapLoop()
    return cache
}

func (c *Cache) Add(key string, val []byte) {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    c.data[key] = cacheEntry{
        createdAt: time.Now(),
        val:       val,
    }
}

func (c *Cache) Get(key string) ([]byte, bool) {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    entry, ok := c.data[key]
    if !ok {
        return nil, false
    }
    return entry.val, true
}

func (c *Cache) reapLoop() {
    ticker := time.NewTicker(c.interval)
    defer ticker.Stop()

    for range ticker.C {
        c.reap()
    }
}

func (c *Cache) reap() {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    now := time.Now()
    for key, entry := range c.data {
        if now.Sub(entry.createdAt) > c.interval {
            delete(c.data, key)
        }
    }
}
