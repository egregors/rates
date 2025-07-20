package rates

import (
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

// Package-level logger for decimal operations
var decimalLogger Logger = log.New(os.Stdout, "[DECIMAL] ", log.LstdFlags)

// SetDecimalLogger sets the logger for decimal operations
func SetDecimalLogger(l Logger) {
	if l != nil {
		decimalLogger = l
	}
}

// Decimal represents a decimal number with arbitrary precision
type Decimal struct {
	value uint64 // scaled by 10^scale
	scale uint8  // decimal places
}

// NewDecimal creates a new Decimal from a float64
func NewDecimal(f float64) Decimal {
	if f == 0 {
		return Decimal{value: 0, scale: 0}
	}

	str := strconv.FormatFloat(f, 'f', -1, 64)
	return NewDecimalFromString(str)
}

// NewDecimalFromString creates a new Decimal from a string
func NewDecimalFromString(s string) Decimal {
	s = strings.TrimSpace(s)
	
	if s == "" {
		decimalLogger.Printf("[WARN] empty string provided to NewDecimalFromString, returning zero")
		return Decimal{value: 0, scale: 0}
	}

	dotIndex := strings.Index(s, ".")
	if dotIndex == -1 {
		val, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			decimalLogger.Printf("[ERROR] failed to parse integer string '%s': %v", s, err)
			return Decimal{value: 0, scale: 0}
		}
		return Decimal{value: val, scale: 0}
	}

	intPart := s[:dotIndex]
	fracPart := s[dotIndex+1:]
	
	intPart = strings.TrimLeft(intPart, "0")
	if intPart == "" {
		intPart = "0"
	}
	
	fracPart = strings.TrimRight(fracPart, "0")
	
	scale := uint8(len(fracPart))
	if scale == 0 {
		val, err := strconv.ParseUint(intPart, 10, 64)
		if err != nil {
			decimalLogger.Printf("[ERROR] failed to parse integer part '%s': %v", intPart, err)
			return Decimal{value: 0, scale: 0}
		}
		return Decimal{value: val, scale: 0}
	}

	combined := intPart + fracPart
	val, err := strconv.ParseUint(combined, 10, 64)
	if err != nil {
		decimalLogger.Printf("[ERROR] failed to parse combined number '%s': %v", combined, err)
		return Decimal{value: 0, scale: 0}
	}

	return Decimal{value: val, scale: scale}
}

// NewDecimalFromParts creates a new Decimal from value and scale
func NewDecimalFromParts(value uint64, scale uint8) Decimal {
	return Decimal{value: value, scale: scale}
}

// Float64 converts the Decimal to a float64
func (d Decimal) Float64() float64 {
	if d.scale == 0 {
		return float64(d.value)
	}
	
	divisor := math.Pow10(int(d.scale))
	return float64(d.value) / divisor
}

// String returns the string representation of the Decimal
func (d Decimal) String() string {
	if d.scale == 0 {
		return strconv.FormatUint(d.value, 10)
	}

	valueStr := strconv.FormatUint(d.value, 10)
	
	for len(valueStr) <= int(d.scale) {
		valueStr = "0" + valueStr
	}
	
	intPart := valueStr[:len(valueStr)-int(d.scale)]
	fracPart := valueStr[len(valueStr)-int(d.scale):]
	
	fracPart = strings.TrimRight(fracPart, "0")
	
	if fracPart == "" {
		return intPart
	}
	
	return intPart + "." + fracPart
}

// Mul multiplies two Decimals
func (d Decimal) Mul(other Decimal) Decimal {
	newValue := d.value * other.value
	newScale := d.scale + other.scale
	
	if newScale > MaxDecimalScale { // Max practical scale
		excess := newScale - MaxDecimalScale
		divisor := uint64(1)
		for i := uint8(0); i < excess; i++ {
			divisor *= 10
		}
		newValue = newValue / divisor
		newScale = MaxDecimalScale
		decimalLogger.Printf("[WARN] precision loss in multiplication, scale reduced from %d to %d", d.scale+other.scale, newScale)
	}
	
	return Decimal{value: newValue, scale: newScale}
}

// Add adds two Decimals
func (d Decimal) Add(other Decimal) Decimal {
	maxScale := d.scale
	if other.scale > maxScale {
		maxScale = other.scale
	}
	
	d1 := d.normalize(maxScale)
	d2 := other.normalize(maxScale)
	
	newValue := d1.value + d2.value
	return Decimal{value: newValue, scale: maxScale}
}

// Sub subtracts two Decimals
func (d Decimal) Sub(other Decimal) Decimal {
	maxScale := d.scale
	if other.scale > maxScale {
		maxScale = other.scale
	}
	
	d1 := d.normalize(maxScale)
	d2 := other.normalize(maxScale)
	
	// Handle negative result (underflow)
	if d1.value < d2.value {
		decimalLogger.Printf("[WARN] subtraction would result in negative value, clamping to zero: %d - %d", d1.value, d2.value)
		return Decimal{value: 0, scale: maxScale}
	}
	
	newValue := d1.value - d2.value
	
	return Decimal{value: newValue, scale: maxScale}
}

// normalize adjusts the decimal to have the specified scale
func (d Decimal) normalize(targetScale uint8) Decimal {
	if d.scale == targetScale {
		return d
	}
	
	if d.scale < targetScale {
		diff := targetScale - d.scale
		multiplier := uint64(math.Pow10(int(diff)))
		return Decimal{value: d.value * multiplier, scale: targetScale}
	}
	
	diff := d.scale - targetScale
	divisor := uint64(math.Pow10(int(diff)))
	return Decimal{value: d.value / divisor, scale: targetScale}
}

// IsZero returns true if the decimal is zero
func (d Decimal) IsZero() bool {
	return d.value == 0
}

// Value returns the internal value (scaled integer)
func (d Decimal) Value() uint64 {
	return d.value
}

// Scale returns the scale (number of decimal places)
func (d Decimal) Scale() uint8 {
	return d.scale
}