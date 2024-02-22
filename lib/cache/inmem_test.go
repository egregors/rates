package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInMem(t *testing.T) {
	t.Parallel()

	// init cache
	c := NewInMem[int](time.Second)
	assert.Equal(t, 0, c.Len())
	assert.Equal(t, map[string]int{}, c.ToMap())

	// try to get non-existing key
	_, ok := c.Get("key")
	assert.False(t, ok)

	// set key
	c.Set("key", 1)
	k, ok := c.Get("key")
	assert.True(t, ok)
	assert.Equal(t, 1, k)
	assert.Equal(t, 1, c.Len())
	assert.Equal(t, map[string]int{"key": 1}, c.ToMap())

	// wait for key to expire
	time.Sleep(time.Second)
	_, ok = c.Get("key")
	assert.False(t, ok)
	assert.Equal(t, 0, c.Len())
	assert.Equal(t, map[string]int{}, c.ToMap())
}
