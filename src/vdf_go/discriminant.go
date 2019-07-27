package vdf_go

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"math/big"
)

type Pair struct {
	p int64
	q int64
}

var odd_primes = primeLessThanN(1 << 16)
var m = 8 * 3 * 5 * 7 * 11 * 13
var residues = make([]int, 0, m)
var sieve_info = make([]Pair, 0, len(odd_primes))

func init() {
	for x := 7; x < m; x += 8 {
		if (x%3 != 0) && (x%5 != 0) && (x%7 != 0) && (x%11 != 0) && (x%13 != 0) {
			residues = append(residues, x)
		}
	}

	var odd_primes_above_13 = odd_primes[5:]

	for i := 0; i < len(odd_primes_above_13); i++ {
		prime := int64(odd_primes_above_13[i])
		sieve_info = append(sieve_info, Pair{p: int64(prime), q: modExp(int64(m)%prime, prime-2, prime)})
	}
}

func modExp(base, exponent, modulus int64) int64 {
	if modulus == 1 {
		return 0
	}
	base = base % modulus
	result := int64(1)
	for i := int64(0); i < exponent; i++ {
		result = (result * base) % modulus
	}
	return result
}

func EntropyFromSeed(seed []byte, byte_count int) []byte {
	buffer := bytes.Buffer{}
	bufferSize := 0

	extra := uint16(0)
	for bufferSize <= byte_count {
		extra_bits := make([]byte, 2)
		binary.BigEndian.PutUint16(extra_bits, extra)
		more_entropy := sha256.Sum256(append(seed, extra_bits...)[:])
		buffer.Write(more_entropy[:])
		bufferSize += sha256.Size
		extra += 1
	}

	return buffer.Bytes()[:byte_count]
}

//Return a discriminant of the given length using the given seed
//It is a random prime p between 13 - 2^2K
//return -p, where p % 8 == 7
func CreateDiscriminant(seed []byte, length int) *big.Int {
	extra := uint8(length) & 7
	byte_count := ((length + 7) >> 3) + 2
	entropy := EntropyFromSeed(seed, byte_count)

	n := new(big.Int)
	n.SetBytes(entropy[:len(entropy)-2])
	n = new(big.Int).Rsh(n, uint(((8 - extra) & 7)))
	n = new(big.Int).SetBit(n, length-1, 1)
	n = new(big.Int).Sub(n, new(big.Int).Mod(n, big.NewInt(int64(m))))
	n = new(big.Int).Add(n, big.NewInt(int64(residues[int(binary.BigEndian.Uint16(entropy[len(entropy)-2:len(entropy)]))%len(residues)])))

	negN := new(big.Int).Neg(n)

	// Find the smallest prime >= n of the form n + m*x
	for {
		sieve := make([]bool, (1 << 16))

		for _, v := range sieve_info {
			// q = m^-1 (mod p)
			// i = -n / m, so that m*i is -n (mod p)
			//i := ((-n % v.p) * v.q) % v.p
			i := (new(big.Int).Mod(negN, big.NewInt(v.p)).Int64() * v.q) % v.p

			for i < int64(len(sieve)) {
				sieve[i] = true
				i += v.p
			}
		}

		for i, v := range sieve {
			t := new(big.Int).Add(n, big.NewInt(int64(m)*int64(i)))
			if !v && t.ProbablyPrime(1) {
				return new(big.Int).Neg(t)
			}
		}

		//n += m * (1 << 16)
		bigM := big.NewInt(int64(m))
		n = new(big.Int).Add(n, bigM.Mul(bigM, big.NewInt(int64(1<<16))))

	}
}
