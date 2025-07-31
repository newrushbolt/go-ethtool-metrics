package generic_info

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPortSpeed(t *testing.T) {
	bits := _GetPortSpeedBits("1Mb/s")
	assert.Equal(t, 1000000.0, bits)

	bits = _GetPortSpeedBits("1Gb/s")
	assert.Equal(t, 1000000000.0, bits)

	bits = _GetPortSpeedBits(" 1Gb/s ")
	assert.Equal(t, 1000000000.0, bits)
}

func TestGetPortSpeedBroken(t *testing.T) {
	bits := _GetPortSpeedBits("Gb/s")
	assert.True(t, math.IsNaN(bits))

	bits = _GetPortSpeedBits("666 WTF/m gg")
	assert.True(t, math.IsNaN(bits))

	bits = _GetPortSpeedBits("6699999999999999999999996 Gb/s")
	assert.True(t, math.IsNaN(bits))
}

func TestEmptyParseInfo(t *testing.T) {
	config := CollectConfig{}.Default()
	result := ParseInfo("", config)
	assert.Nil(t, result)
}
