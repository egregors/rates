package rates

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestWithLogger(t *testing.T) {
	// if WithLogger is not called, the default logger is used
	assert.NotNil(t, New(nil).l)

	// if WithLogger is called, the logger is used
	logger := &log.Logger{}
	assert.Equal(t, logger, New(nil, WithLogger(logger)).l)

	// if WithLogger is called with nil, the default logger is used
	assert.NotNil(t, New(nil, WithLogger(nil)).l)
}
