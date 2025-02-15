package cache

import (
	"sync"
	"time"
)

type InMem[T any] struct {
	cache   map[string]T
	ttl     time.Duration
	setTime map[string]time.Time

	mu sync.RWMutex
}

func NewInMem[T any](ttl time.Duration) *InMem[T] {
	return &InMem[T]{
		cache:   make(map[string]T),
		ttl:     ttl,
		setTime: make(map[string]time.Time),
		mu:      sync.RWMutex{},
	}
}

func (c *InMem[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if t, ok := c.setTime[key]; ok {
		if time.Since(t) > c.ttl {
			delete(c.cache, key)
			delete(c.setTime, key)

			return *new(T), false
		}
	}

	val, ok := c.cache[key]

	return val, ok
}

func (c *InMem[T]) Set(key string, val T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = val
	c.setTime[key] = time.Now()
}

func (c *InMem[T]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.cache)
}

func (c *InMem[T]) ToMap() map[string]T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.cache
}
