package rates

import (
	"github.com/egregors/rates/lib"
	"github.com/egregors/rates/lib/cache"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithLogger(t *testing.T) {
	// if WithLogger is not called, the default logger is used
	assert.NotNil(t, New(nil).l)
	assert.IsType(t, &log.Logger{}, New(nil).l)

	// if WithLogger is called, the logger is used
	logger := &log.Logger{}
	assert.Equal(t, logger, New(nil, WithLogger(logger)).l)
	assert.IsType(t, &log.Logger{}, New(nil, WithLogger(logger)).l)

	// if WithLogger is called with nil, the default logger is used
	assert.NotNil(t, New(nil, WithLogger(nil)).l)
	assert.IsType(t, &lib.NoopLogger{}, New(nil, WithLogger(nil)).l)
}

func TestWithCache(t *testing.T) {
	// if WithCache is not called, the default cache is used
	assert.NotNil(t, New(nil).cache)
	assert.IsType(t, &cache.InMem[map[string]float64]{}, New(nil).cache)

	//if WithCache is called, the cache is used
	c := cache.NewInMem[map[string]float64](69)
	assert.Equal(t, c, New(nil, WithCache(c)).cache)
	assert.IsType(t, &cache.InMem[map[string]float64]{}, New(nil, WithCache(c)).cache)

	// if WithCache is called with nil, the default cache is used
	assert.NotNil(t, New(nil, WithCache(nil)).cache)
	assert.IsType(t, &cache.Noop[map[string]float64]{}, New(nil, WithCache(nil)).cache)
}
