package vacuum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinUpdate(t *testing.T) {
	assert := assert.New(t)

	b := Bin{
		FilePath: `../tests/RoboController.cfg`,
		Capacity: 3600.,
	}

	b.Update()
	assert.Equal(b.Seconds, 100.)
}

func TestConvert(t *testing.T) {
	assert := assert.New(t)
	var tests = []struct {
		unit     string
		capacity float64
		expected string
	}{
		{unit: "%", capacity: 3600., expected: "2.78"},
		{unit: "%", capacity: 2400., expected: "4.17"},
		{unit: "sec", capacity: 3600., expected: "100"},
		{unit: "min", capacity: 3600., expected: "2"},
	}

	for _, test := range tests {
		b := Bin{
			FilePath: `../tests/RoboController.cfg`,
			Capacity: test.capacity,
			Unit:     test.unit,
		}

		b.Update()
		assert.Equal(100., b.Seconds)
		b.convert()
		assert.Equal(test.expected, b.Value)

	}

}
