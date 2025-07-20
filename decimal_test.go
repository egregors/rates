package rates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDecimal(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{"zero", 0.0, "0"},
		{"integer", 123.0, "123"},
		{"decimal", 123.456, "123.456"},
		{"small decimal", 0.123, "0.123"},
		{"large number", 1234567.89, "1234567.89"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDecimal(tt.input)
			assert.Equal(t, tt.expected, d.String())
		})
	}
}

func TestNewDecimalFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"zero", "0", "0"},
		{"integer", "123", "123"},
		{"decimal", "123.456", "123.456"},
		{"small decimal", "0.123", "0.123"},
		{"trailing zeros", "123.4500", "123.45"},
		{"leading zeros", "0123.456", "123.456"},
		{"empty string", "", "0"},
		{"whitespace", "  123.45  ", "123.45"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDecimalFromString(tt.input)
			assert.Equal(t, tt.expected, d.String())
		})
	}
}

func TestNewDecimalFromParts(t *testing.T) {
	tests := []struct {
		name     string
		value    uint64
		scale    uint8
		expected string
	}{
		{"zero", 0, 0, "0"},
		{"integer", 123, 0, "123"},
		{"decimal", 12345, 3, "12.345"},
		{"small", 123, 5, "0.00123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDecimalFromParts(tt.value, tt.scale)
			assert.Equal(t, tt.expected, d.String())
		})
	}
}

func TestFloat64(t *testing.T) {
	tests := []struct {
		name     string
		decimal  Decimal
		expected float64
	}{
		{"zero", NewDecimal(0), 0.0},
		{"integer", NewDecimal(123), 123.0},
		{"decimal", NewDecimal(123.456), 123.456},
		{"small decimal", NewDecimal(0.123), 0.123},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.decimal.Float64()
			assert.InDelta(t, tt.expected, result, 0.000001)
		})
	}
}

func TestDecimalMul(t *testing.T) {
	tests := []struct {
		name     string
		d1       Decimal
		d2       Decimal
		expected string
	}{
		{"zero", NewDecimal(0), NewDecimal(123), "0"},
		{"integer", NewDecimal(10), NewDecimal(5), "50"},
		{"decimal", NewDecimal(12.5), NewDecimal(2.4), "30"},
		{"precision", NewDecimal(1.23), NewDecimal(4.56), "5.6088"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.d1.Mul(tt.d2)
			assert.Equal(t, tt.expected, result.String())
		})
	}
}

func TestDecimalAdd(t *testing.T) {
	tests := []struct {
		name     string
		d1       Decimal
		d2       Decimal
		expected string
	}{
		{"zero", NewDecimal(0), NewDecimal(123), "123"},
		{"integer", NewDecimal(10), NewDecimal(5), "15"},
		{"decimal", NewDecimal(12.5), NewDecimal(2.4), "14.9"},
		{"different scales", NewDecimal(1.2), NewDecimal(3.456), "4.656"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.d1.Add(tt.d2)
			assert.Equal(t, tt.expected, result.String())
		})
	}
}

func TestDecimalSub(t *testing.T) {
	tests := []struct {
		name     string
		d1       Decimal
		d2       Decimal
		expected string
	}{
		{"zero", NewDecimal(123), NewDecimal(0), "123"},
		{"integer", NewDecimal(10), NewDecimal(5), "5"},
		{"decimal", NewDecimal(12.5), NewDecimal(2.4), "10.1"},
		{"different scales", NewDecimal(3.456), NewDecimal(1.2), "2.256"},
		{"negative result", NewDecimal(1), NewDecimal(2), "0"}, // Clamps to zero
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.d1.Sub(tt.d2)
			assert.Equal(t, tt.expected, result.String())
		})
	}
}

func TestDecimalIsZero(t *testing.T) {
	assert.True(t, NewDecimal(0).IsZero())
	assert.False(t, NewDecimal(1).IsZero())
	assert.False(t, NewDecimal(0.001).IsZero())
}

func TestDecimalValueScale(t *testing.T) {
	d := NewDecimalFromParts(12345, 3)
	assert.Equal(t, uint64(12345), d.Value())
	assert.Equal(t, uint8(3), d.Scale())
}