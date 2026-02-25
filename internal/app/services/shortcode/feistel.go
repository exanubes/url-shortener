package shortcode

import "math/big"

// Constants from SplitMix64 (Sebastiano Vigna)
// Used for key expansion.
const (
	split_mix_64_gamma = 0x9e3779b97f4a7c15
	split_mix_64_mix1  = 0xbf58476d1ce4e5b9
	split_mix_64_mix2  = 0x94d049bb133111eb
)

// Constants from MurmurHash3 64-bit finalizer.
// Used in round function for strong avalanche.
const (
	murmur3_finalizer_mix1 = 0xff51afd7ed558ccd
	murmur3_finalizer_mix2 = 0xc4ceb9fe1a85ec53
)

type Feistel struct {
	roundKeys [8]uint64
	rounds    uint8
}

// Create a new Feistel permutation with a 64-bit master key
func NewFeistel(key uint64) *Feistel {
	var f Feistel
	f.rounds = 8
	// Derive round keys using splitmix64
	k := key
	for i := range f.rounds {
		k += split_mix_64_gamma
		z := k
		z = (z ^ (z >> 30)) * split_mix_64_mix1
		z = (z ^ (z >> 27)) * split_mix_64_mix2
		z = z ^ (z >> 31)
		f.roundKeys[i] = z
	}

	return &f
}

func round_function(r uint32, k uint64) uint32 {
	x := uint64(r) ^ k
	x ^= x >> 33
	x *= murmur3_finalizer_mix1
	x ^= x >> 33
	x *= murmur3_finalizer_mix2
	x ^= x >> 33
	return uint32(x)
}

func (f *Feistel) Scramble(input Token) (Token, error) {
	scrambled := f.Encrypt(input.value.Uint64())
	bigint := (new(big.Int)).SetUint64(scrambled)
	return NewToken(bigint, 11)
}

// Encrypt (permute)
func (f *Feistel) Encrypt(v uint64) uint64 {
	l := uint32(v >> 32)
	r := uint32(v)

	for i := range f.rounds {
		newL := r
		newR := l ^ round_function(r, f.roundKeys[i])
		l = newL
		r = newR
	}

	return (uint64(l) << 32) | uint64(r)
}

// Decrypt (reverse permutation)
func (f *Feistel) Decrypt(v uint64) uint64 {
	l := uint32(v >> 32)
	r := uint32(v)

	for i := 7; i >= 0; i-- {
		newR := l
		newL := r ^ round_function(l, f.roundKeys[i])
		l = newL
		r = newR
	}

	return (uint64(l) << 32) | uint64(r)
}
