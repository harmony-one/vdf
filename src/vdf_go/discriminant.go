package vdf_go

import (
	"crypto/sha256"
	"encoding/binary"
	"math/big"
)

type Pair struct {
	p int64
	q int64
}

var m = 8 * 3 * 5 * 7 * 11 * 13

var (
	bigMValue     = big.NewInt(int64(m))
	bigSieveStep  = new(big.Int).Lsh(new(big.Int).Set(bigMValue), 16)
	sieveInfoBigP = buildSieveInfoBigP()
)

func buildSieveInfoBigP() []*big.Int {
	values := make([]*big.Int, len(sieve_info))
	for i, pair := range sieve_info {
		values[i] = big.NewInt(pair.p)
	}
	return values
}

func EntropyFromSeed(seed []byte, byte_count int) []byte {
	extra := uint16(0)
	input := make([]byte, len(seed)+2)
	copy(input, seed)
	output := make([]byte, byte_count+sha256.Size)

	offset := 0
	for offset <= byte_count {
		binary.BigEndian.PutUint16(input[len(seed):], extra)
		sum := sha256.Sum256(input)
		copy(output[offset:], sum[:])
		offset += sha256.Size
		extra++
	}

	return output[:byte_count]
}

// Return a discriminant of the given length using the given seed
// It is a random prime p between 13 - 2^2K
// return -p, where p % 8 == 7
func CreateDiscriminant(seed []byte, length int) *big.Int {
	extra := uint8(length) & 7
	byte_count := ((length + 7) >> 3) + 2
	entropy := EntropyFromSeed(seed, byte_count)

	n := new(big.Int).SetBytes(entropy[:len(entropy)-2])
	n.Rsh(n, uint((8-extra)&7))
	n.SetBit(n, length-1, 1)

	mod := new(big.Int).Mod(n, bigMValue)
	n.Sub(n, mod)
	n.Add(n, big.NewInt(int64(residues[int(binary.BigEndian.Uint16(entropy[len(entropy)-2:]))%len(residues)])))

	negN := new(big.Int).Neg(n)
	modP := new(big.Int)
	sieve := make([]bool, 1<<16)

	// Find the smallest prime >= n of the form n + m*x
	for {
		clear(sieve)

		for idx, v := range sieve_info {
			// q = m^-1 (mod p)
			// i = -n / m, so that m*i is -n (mod p)
			//i := ((-n % v.p) * v.q) % v.p
			i := (modP.Mod(negN, sieveInfoBigP[idx]).Int64() * v.q) % v.p

			for i < int64(len(sieve)) {
				sieve[i] = true
				i += v.p
			}
		}

		t := new(big.Int).Set(n)
		for _, composite := range sieve {
			if !composite && t.ProbablyPrime(1) {
				return new(big.Int).Neg(t)
			}
			t.Add(t, bigMValue)
		}

		//n += m * (1 << 16)
		n.Add(n, bigSieveStep)
		negN.Neg(n)
	}
}
