package cachey

import (
	"bytes"
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	c := NewCache()
	c.Set("key1", "val1", time.Second)
	c.Set("key2", "val2", 2*time.Second)
	c.Set("key3", "val3", time.Second)

	if v, ok := c.Get("key1"); !ok || v != "val1" {
		t.Fatal("not set")
	}

	time.Sleep(1 * time.Second)

	if _, ok := c.Get("key1"); ok {
		t.Fatal("not expired!")
	}
	if _, ok := c.Get("key2"); !ok {
		t.Fatal("expired!")
	}
	if _, ok := c.Get("key3"); ok {
		t.Fatal("not expired!")
	}
}

func TestValueIsNotFound(t *testing.T) {
	c := NewCache()
	if v, ok := c.Get("not-found"); v != nil || ok {
		t.Fatal("found 'not-found-key's cache.")
	}
}

func TestGetOrSet(t *testing.T) {
	c := NewCache()
	key := "get_or_set"
	if v, ok := c.Get(key); v != nil || ok {
		t.Fatal("found 'get_or_set's cache.")
	}

	value, isGot := c.GetOrSet(key, func() (interface{}, time.Duration) {
		return "value1", time.Second
	})
	if value != "value1" || isGot {
		t.Fatal("not set")
	}

	value, isGot = c.GetOrSet(key, func() (interface{}, time.Duration) {
		return "value2", time.Second
	})
	if value == "value2" || !isGot {
		t.Fatal("set value2")
	}

	time.Sleep(1 * time.Second)

	// expired. so, not found cache.
	if value, ok := c.Get(key); value != nil || ok {
		t.Fatal("not expired")
	}
}

func TestDelete(t *testing.T) {
	c := NewCache()
	value, _ := c.GetOrSet("name", func() (interface{}, time.Duration) {
		return "mattn", 10 * time.Second
	})
	if value == "" {
		t.Fatal("not set")
	}

	if deletedValue, ok := c.Delete("name"); deletedValue == "" || !ok {
		t.Fatal("not delete")
	}
}

func TestLock(t *testing.T) {
	c := NewCache()
	c.Set("name", "mattn", 10*time.Second)

	ch := make(chan string)
	go func() {
		c.Get("name")
		ch <- "1"
	}()
	go func() {
		c.Get("name")
		ch <- "2"
	}()
	go func() {
		c.Set("name", "mattn", 10*time.Second)
		ch <- "3"
	}()
	go func() {
		c.Get("name")
		ch <- "4"
	}()
	go func() {
		c.Delete("name")
		ch <- "5"
	}()
	go func() {
		c.Get("name")
		ch <- "6"
	}()

	var buffer bytes.Buffer
	for i := 0; i < 6; i++ {
		buffer.WriteString(<-ch)
	}
	if buffer.String() != "123456" {
		t.Fatal("not locked", buffer.String())
	}
}

// This is based on https://github.com/hisaichi5518/vache
