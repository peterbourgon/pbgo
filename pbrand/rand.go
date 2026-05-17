package pbrand

import (
	"crypto/sha256"
	"encoding"
	"encoding/binary"
	"io"
	"math/rand/v2"
)

type ChaCha8Rand struct {
	src *rand.ChaCha8
	*rand.Rand
}

var (
	_ rand.Source                    = (*ChaCha8Rand)(nil)
	_ io.Reader                      = (*ChaCha8Rand)(nil)
	_ encoding.BinaryUnmarshaler     = (*ChaCha8Rand)(nil)
	_ interface{ IntN(int) int }     = (*ChaCha8Rand)(nil)
	_ interface{ Float64() float64 } = (*ChaCha8Rand)(nil)
)

func NewChaCha8Rand(seed int64) *ChaCha8Rand {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(seed))
	hash := sha256.Sum256(buf[:])
	src := rand.NewChaCha8(hash)
	rand := rand.New(src)
	return &ChaCha8Rand{
		src:  src,
		Rand: rand,
	}
}

func (r *ChaCha8Rand) Uint64() uint64 {
	return r.src.Uint64()
}

func (r *ChaCha8Rand) Read(p []byte) (n int, err error) {
	return r.src.Read(p)
}

func (r *ChaCha8Rand) UnmarshalBinary(data []byte) error {
	return r.src.UnmarshalBinary(data)
}
