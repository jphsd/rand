package rand

import (
	"math/bits"
	"math/rand"
	"sync"
)

const (
	// MaxUint63 represents the largest value that can be held in a 63 bit unsigned integer.
	MaxUint63 = (1 << 63) - 1

	mulHi = 0x2360ed051fc65da4
	mulLo = 0x4385df649fccf645
	incHi = 0x5851f42d4c957f2d
	incLo = 0x14057b7ef767814f
)

// SplitMix64Source - size 64 bits, period 2^64.
type SplitMix64Source uint64

// Seed implements the rand.Source interface.
func (s *SplitMix64Source) Seed(seed int64) {
	*s = SplitMix64Source(seed)
}

// Int63 implements the rand.Source interface.
func (s *SplitMix64Source) Int63() int64 {
	return int64(s.Uint64() & MaxUint63)
}

// Uint64 implements the rand.Source64 interface.
func (s *SplitMix64Source) Uint64() uint64 {
	*s += 0x9e3779b97f4a7c15
	z := *s
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5B9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eB
	z = z ^ (z >> 31)
	return uint64(z)
}

// XOshiro256Source - size 256 bits, period 2^256 - 1.
// Do not leave zero initialized!
type XOshiro256Source [4]uint64

// Seed implements the rand.Source interface.
func (s *XOshiro256Source) Seed(seed int64) {
	var sm64 SplitMix64Source
	sm64.Seed(seed)
	s[0], s[1], s[2], s[3] = sm64.Uint64(), sm64.Uint64(), sm64.Uint64(), sm64.Uint64()
}

// Int63 implements the rand.Source interface.
func (s *XOshiro256Source) Int63() int64 {
	return int64(s.Uint64() & MaxUint63)
}

// Uint64 implements the rand.Source64 interface.
func (s *XOshiro256Source) Uint64() uint64 {
	z := bits.RotateLeft64(s[1]*5, 7) * 9
	t := s[1] << 17

	s[2] ^= s[0]
	s[3] ^= s[1]
	s[1] ^= s[2]
	s[0] ^= s[3]

	s[2] ^= t
	s[3] = bits.RotateLeft64(s[3], 45)

	return z
}

// Jump is equivalent to 2^128 calls to Uint64(); it can be used to generate 2^128
// non-overlapping subsequences for parallel computations.
func (s *XOshiro256Source) Jump() *XOshiro256Source {
	jbits := [4]uint64{0x180ec6d33cfd0aba, 0xd5a61266f0c9392c, 0xa9582618e03fc9aa, 0x39abdc4529b1661c}

	var s0, s1, s2, s3 uint64
	for i := 0; i < len(jbits); i++ {
		for b := 0; b < 64; b++ {
			if jbits[i]&(1<<b) != 0 {
				s0 ^= s[0]
				s1 ^= s[1]
				s2 ^= s[2]
				s3 ^= s[3]
			}
			s.Uint64()
		}
	}

	return &XOshiro256Source{s0, s1, s2, s3}
}

// LongJump is equivalent to 2^192 calls to Uint64(); it can be used to generate 2^64 starting points,
// from each of which jump() will generate 2^64 non-overlapping subsequences for parallel distributed
// computations.
func (s *XOshiro256Source) LongJump() *XOshiro256Source {
	jbits := [4]uint64{0x76e15d3efefdcbbf, 0xc5004e441c522fb3, 0x77710069854ee241, 0x39109bb02acbe635}

	var s0, s1, s2, s3 uint64
	for i := 0; i < len(jbits); i++ {
		for b := 0; b < 64; b++ {
			if jbits[i]&(1<<b) != 0 {
				s0 ^= s[0]
				s1 ^= s[1]
				s2 ^= s[2]
				s3 ^= s[3]
			}
			s.Uint64()
		}
	}

	return &XOshiro256Source{s0, s1, s2, s3}
}

// PCGSource - size 128 bits, period 2^128.
type PCGSource [2]uint64

// Seed implements the rand.Source interface.
func (s *PCGSource) Seed(seed int64) {
	var sm64 SplitMix64Source
	sm64.Seed(seed)
	s[0], s[1] = sm64.Uint64(), sm64.Uint64()
}

// Int63 implements the rand.Source interface.
func (s *PCGSource) Int63() int64 {
	return int64(s.Uint64() & MaxUint63)
}

// Uint64 implements the rand.Source64 interface.
func (s *PCGSource) Uint64() uint64 {
	s.mult()
	s.add()
	return bits.RotateLeft64(s[0]^s[1], -int(s[1]>>58))
}

func (s *PCGSource) add() {
	var c uint64
	s[0], c = bits.Add64(s[0], incLo, 0)
	s[1], _ = bits.Add64(s[1], incHi, c)
}

func (s *PCGSource) mult() {
	hi, lo := bits.Mul64(s[0], mulLo)
	hi += s[1] * mulLo
	hi += s[0] * mulHi
	s[0], s[1] = lo, hi
}

// wyrand (runtime.fastrand)

const (
	wy0 WySource = 0xa0761d6478bd642f
	wy1 WySource = 0xe7037ed1a0b428db
)

// WySource (from https://github.com/wangyi-fudan/wyhash)
type WySource uint64

// Seed implements the rand.Source interface.
func (s *WySource) Seed(seed int64) {
	*s = WySource(seed)
}

// Int63 implements the rand.Source interface.
func (s *WySource) Int63() int64 {
	return int64(s.Uint64() & MaxUint63)
}

// Uint64 implements the rand.Source64 interface.
func (s *WySource) Uint64() uint64 {
	*s += wy0
	return wymulf(*s^wy1, *s)
}

func wymulf(a, b WySource) uint64 {
	hi, lo := bits.Mul64(uint64(a), uint64(b))
	return hi ^ lo
}

// LockableSource wraps a type implementing rand.Source64 in a mutex to make it goroutine safe.
type LockableSource struct {
	lk  sync.Mutex
	src rand.Source64
}

// NewLockableSource returns src wrapped in a mutex.
func NewLockableSource(src rand.Source64) *LockableSource {
	return &LockableSource{src: src}
}

// Seed implements the rand.Source interface.
func (s *LockableSource) Seed(seed int64) {
	s.lk.Lock()
	s.src.Seed(seed)
	s.lk.Unlock()
}

// Int63 implements the rand.Source interface.
func (s *LockableSource) Int63() int64 {
	s.lk.Lock()
	z := s.src.Int63()
	s.lk.Unlock()
	return z
}

// Uint64 implements the rand.Source64 interface.
func (s *LockableSource) Uint64() uint64 {
	s.lk.Lock()
	z := s.src.Uint64()
	s.lk.Unlock()
	return z
}
