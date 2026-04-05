package main

import (
	"testing"

	"github.com/harmony-one/vdf/src/vdf_go"
	"github.com/stretchr/testify/assert"
)

func buildValidProofFixture(t *testing.T) (seed []byte, proofBlob []byte, iterations int, intSizeBits int) {
	seed = []byte{
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
	}
	iterations = 20
	intSizeBits = 2048

	yBuf, proofBuf := vdf_go.GenerateVDF(seed, iterations, intSizeBits)
	assert.NotNil(t, yBuf)
	assert.NotNil(t, proofBuf)

	proofBlob = append(append([]byte{}, yBuf...), proofBuf...)
	assert.Equal(t, true, vdf_go.VerifyVDF(seed, proofBlob, iterations, intSizeBits), "fixture proof should be valid")

	return seed, proofBlob, iterations, intSizeBits
}

func TestVerifyProofTamperedBlob(t *testing.T) {
	seed, proofBlob, iterations, intSizeBits := buildValidProofFixture(t)

	tampered := append([]byte{}, proofBlob...)
	tampered[len(tampered)-1] ^= 0x01

	assert.False(t, vdf_go.VerifyVDF(seed, tampered, iterations, intSizeBits), "tampered proof must fail")
}

func TestVerifyProofWrongTarget(t *testing.T) {
	seed, proofBlob, iterations, intSizeBits := buildValidProofFixture(t)

	assert.False(t, vdf_go.VerifyVDF(seed, proofBlob, iterations+1, intSizeBits), "wrong iteration target must fail")
}

func TestVerifyProofWrongSeed(t *testing.T) {
	seed, proofBlob, iterations, intSizeBits := buildValidProofFixture(t)

	wrongSeed := append([]byte{}, seed...)
	wrongSeed[0] ^= 0x01

	assert.False(t, vdf_go.VerifyVDF(wrongSeed, proofBlob, iterations, intSizeBits), "wrong seed must fail")
}

func TestVerifyProofWrongIntSizeBits(t *testing.T) {
	seed, proofBlob, iterations, intSizeBits := buildValidProofFixture(t)

	assert.Panics(t, func() {
		vdf_go.VerifyVDF(seed, proofBlob, iterations, intSizeBits-1)
	}, "wrong int_size_bits currently panics with current API contract")
}

func TestVerifyProofShortBlobPanics(t *testing.T) {
	seed, _, iterations, intSizeBits := buildValidProofFixture(t)

	shortBlob := []byte{0x01, 0x02, 0x03}

	assert.Panics(t, func() {
		vdf_go.VerifyVDF(seed, shortBlob, iterations, intSizeBits)
	}, "too short proof blob should panic with current API contract")
}
