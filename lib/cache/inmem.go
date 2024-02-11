package cache

import (
	"sync"
)

type InMem[T any] struct {
	cache map[string]T
	mu    sync.RWMutex
}

func NewInMem[T any]() *InMem[T] {
	return &InMem[T]{
		cache: make(map[string]T),
		mu:    sync.RWMutex{},
	}
}

func (c *InMem[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, ok := c.cache[key]

	return val, ok
}

func (c *InMem[T]) Set(key string, val T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = val
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
