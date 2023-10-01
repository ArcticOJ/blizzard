package numeric

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestCompress(t *testing.T) {
	for i := uint16(0); i < math.MaxUint16; i++ {
		for j := uint16(0); j < math.MaxUint16; j++ {
			go func(a, b uint16) {
				compressed := CompressUint16(i, j)
				_a, _b := DecompressUint32(compressed)
				assert.Truef(t, a == _a && b == _b, "expected true with a = %d and b = %d, got a = %d, b = %d", a, b, _a, _b)
			}(i, j)
		}
	}
}
