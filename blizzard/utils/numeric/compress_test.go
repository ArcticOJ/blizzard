package numeric

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestCompress(t *testing.T) {
	for i := uint16(0); i < math.MaxUint16; i++ {
		for j := uint16(0); j < math.MaxUint16; j++ {
			compressed := CompressUint16(i, j)
			a, b := DecompressUint32(compressed)
			assert.Truef(t, i == a && j == b, "expected true with a = %d and b = %d, got a = %d, b = %d", i, j, a, b)
		}
	}
}
