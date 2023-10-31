package numeric

import (
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"slices"
	"testing"
	"time"
)

type tuple struct {
	first  uint16
	second uint16
	third  uint32
}

func TestCompress(t *testing.T) {
	a1 := make([]tuple, 3e7)
	a2 := make([]tuple, 3e7)
	for i := range a1 {
		p := tuple{first: uint16(rand.Intn(math.MaxUint16)), second: uint16(rand.Intn(math.MaxUint16))}
		p.third = CompressUint16(p.first, p.second)
		a1[i] = p
		a2[i] = p
	}
	slices.SortFunc(a1, func(a, b tuple) int {
		if a.first == b.first {
			return int(a.second - b.second)
		}
		return int(a.first - b.first)
	})
	slices.SortFunc(a2, func(a, b tuple) int {
		return int(a.third - b.third)
	})
	assert.True(t, slices.Equal(a1, a2), "expected a1 to be equal to a2")
}

func BenchmarkDecompress(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < b.N; i++ {
		DecompressUint32(r.Uint32())
	}
}

func BenchmarkCompress(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < b.N; i++ {
		u := uint16(r.Intn(math.MaxUint16))
		CompressUint16(u, u)
	}
}
