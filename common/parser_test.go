package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFloat64(t *testing.T) {
	floatString := "3.14159265359"
	float := parseFloat64(floatString)
	// float64 in Go seems to be precise enough for comparing with Equal()
	assert.Equal(t, 3.14159265359, float)

	floatString = "3"
	float = parseFloat64(floatString)
	assert.Equal(t, 3.0, float)

	floatString = "+3"
	float = parseFloat64(floatString)
	assert.Equal(t, 3.0, float)

	floatString = "-3"
	float = parseFloat64(floatString)
	assert.Equal(t, -3.0, float)

	floatString = "some text 3.14 wont harm"
	float = parseFloat64(floatString)
	assert.Equal(t, 3.14, float)

	floatString = "nice text wont float"
	float = parseFloat64(floatString)
	assert.Equal(t, 0.0, float)

	floatString = "3.14 two floats / 3.15 one string"
	float = parseFloat64(floatString)
	assert.Equal(t, 3.14, float)
}

func TestParseBoolTrue(t *testing.T) {
	boolIsh := "On"
	boolVal := parseBool(boolIsh)
	assert.True(t, boolVal)

	boolIsh = "ON"
	boolVal = parseBool(boolIsh)
	assert.True(t, boolVal)

	boolIsh = " on "
	boolVal = parseBool(boolIsh)
	assert.True(t, boolVal)

	boolIsh = "TRUE"
	boolVal = parseBool(boolIsh)
	assert.True(t, boolVal)

	boolIsh = "yes"
	boolVal = parseBool(boolIsh)
	assert.True(t, boolVal)
}

func TestParseBoolFalse(t *testing.T) {
	boolIsh := "no"
	boolVal := parseBool(boolIsh)
	assert.False(t, boolVal)

	boolIsh = "NO"
	boolVal = parseBool(boolIsh)
	assert.False(t, boolVal)

	boolIsh = "off"
	boolVal = parseBool(boolIsh)
	assert.False(t, boolVal)

	boolIsh = "oFF"
	boolVal = parseBool(boolIsh)
	assert.False(t, boolVal)

	boolIsh = "false"
	boolVal = parseBool(boolIsh)
	assert.False(t, boolVal)

	// Shouldn't be parsed
	boolIsh = " ON OFF"
	boolVal = parseBool(boolIsh)
	assert.False(t, boolVal)
}

func TestParseSlice(t *testing.T) {
	expectedSlice := []string{"mode1", "mode2", "mode3"}
	sliceString := "mode1 mode2 mode3"
	parsedSlice := parseSlice(sliceString)
	assert.Equal(t, expectedSlice, parsedSlice)

	sliceString = "mode1, mode2, mode3"
	parsedSlice = parseSlice(sliceString)
	assert.Equal(t, expectedSlice, parsedSlice)

	sliceString = "mode1,mode2,mode3"
	parsedSlice = parseSlice(sliceString)
	assert.Equal(t, expectedSlice, parsedSlice)

	sliceString = "mode1\nmode2\tmode3"
	parsedSlice = parseSlice(sliceString)
	assert.Equal(t, expectedSlice, parsedSlice)
}
