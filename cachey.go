package cachey

import (
	"runtime"
	"sync"
	"time"
)

type item struct {
	key    string
	value  interface{}
	expire time.Time
}

type cache struct {
	items        map[string]item
	mutex        sync.RWMutex
	timer        *time.Timer
	DurationOfGC time.Duration
}

// NewCache returns a new cache. After calling NewCache, GC is started to
// expire old items
func NewCache() *cache {
	c := &cache{}
	c.items = make(map[string]item)
	c.DurationOfGC = 5
	c.startGC()
	return c
}

// Set store the value with specified key. This value is stored until while
// expire's seconds.
func (c *cache) Set(key string, value interface{}, expire time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.items[key] = item{
		value:  value,
		expire: time.Now().Add(expire),
	}
}

// Get returns the value that is stored by specified key, and whether the value
// is stored or not. If the value is not flushed, expire time will be updated
// to next duration.
func (c *cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if item.expire.Before(time.Now()) {
		delete(c.items, key)
		return nil, false
	}
	return item.value, true

}

// GetOrSet try to get and set combination.
func (c *cache) GetOrSet(key string, code func() (interface{}, time.Duration)) (interface{}, bool) {
	if value, ok := c.Get(key); ok {
		return value, true
	}
	value, expire := code()
	c.Set(key, value, expire)
	return value, false
}

// Delete deletes the value specified by key.
func (c *cache) Delete(key string) (deleted interface{}, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	deleted, ok = c.items[key]
	delete(c.items, key)
	return
}

// startGC make a timer of GC. It is possible to specify the duration of GC by
// DurationOfGC.
func (c *cache) startGC() {
	var t *time.Timer
	runtime.SetFinalizer(c, func(c *cache) {
		t.Stop()
	})
	t = time.AfterFunc(c.DurationOfGC, func() {
		keys := []string{}
		c.mutex.RLock()
		for k := range c.items {
			keys = append(keys, k)
		}
		c.mutex.RUnlock()

		now := time.Now()

		for _, k := range keys {
			c.mutex.Lock()
			if v, ok := c.items[k]; ok {
				if v.expire.Before(now) {
					delete(c.items, k)
				}
			}
			c.mutex.Unlock()
		}
		t.Reset(c.DurationOfGC)
	})
}
