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

func NewCache() *cache {
	c := &cache{}
	c.items = make(map[string]item)
	c.DurationOfGC = 5
	c.startGC()
	return c
}

func (c *cache) Set(key string, value interface{}, expire time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.items[key] = item{
		value:  value,
		expire: time.Now().Add(expire),
	}
}

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

func (c *cache) GetOrSet(key string, code func() (interface{}, time.Duration)) (interface{}, bool) {
	if value, ok := c.Get(key); ok {
		return value, true
	}
	value, expire := code()
	c.Set(key, value, expire)
	return value, false
}

func (c *cache) Delete(key string) (deleted interface{}, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	deleted, ok = c.items[key]
	delete(c.items, key)
	return
}

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
			v := c.items[k]
			if v.expire.Before(now) {
				delete(c.items, k)
			}
			c.mutex.Unlock()
		}
		t.Reset(c.DurationOfGC)
	})
}
