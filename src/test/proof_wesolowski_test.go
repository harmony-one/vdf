package main

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
	"vdf_go"
)

func TestCreateProof(t *testing.T) {
	seed := []byte{0xde, 0xad, 0xbe, 0xef}

	y_buf, proof_buf := vdf_go.GenerateVDF(seed, 10, 2048)

	str_y := hex.EncodeToString(y_buf)
	str_proof := hex.EncodeToString(proof_buf)

	assert.Equal(t, "0038c64a357958a41b5a3374c0373a93969766295b3e743af7501397c4add9e8f09c40359f2ae2621ba33a85dacf300dc241c27516248a27d562e09000fe5f6f20d0cd8fb05c77bc01e750d1ff3bdec7ce1a9cbc4c8bcf215c4a96bdef9a386b07d7a93f79574740243bfcce2e13c7bb959c36b7077bbcee5e8100a7f8e06bf2d80012a7c998dbff6cb4bbc6cbb36dd7a8b04d018e5870ccfe8416ce6507284697b0082dd2450ad170e41d9f038e5c35ec88501f7760455c875cde50a0609ed3b2f13a6203d1eac326e64e35b95f40a004ca17ed3a59d048dbbb16de24814d4b9462aa91a41d803815fba7e4186e6712bf77d4324ab692c3db4df82c670f49409f3b", str_y, "y match")
	assert.Equal(t, "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001", str_proof, "proof match")
}

func TestVerifyProof(t *testing.T) {
	seed := []byte{0xde, 0xad, 0xbe, 0xef}

	y_buf, _ := hex.DecodeString("0038c64a357958a41b5a3374c0373a93969766295b3e743af7501397c4add9e8f09c40359f2ae2621ba33a85dacf300dc241c27516248a27d562e09000fe5f6f20d0cd8fb05c77bc01e750d1ff3bdec7ce1a9cbc4c8bcf215c4a96bdef9a386b07d7a93f79574740243bfcce2e13c7bb959c36b7077bbcee5e8100a7f8e06bf2d80012a7c998dbff6cb4bbc6cbb36dd7a8b04d018e5870ccfe8416ce6507284697b0082dd2450ad170e41d9f038e5c35ec88501f7760455c875cde50a0609ed3b2f13a6203d1eac326e64e35b95f40a004ca17ed3a59d048dbbb16de24814d4b9462aa91a41d803815fba7e4186e6712bf77d4324ab692c3db4df82c670f49409f3b")
	proof_buf, _ := hex.DecodeString("000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001")

	assert.Equal(t, true, vdf_go.VerifyVDF(seed, append(y_buf, proof_buf...), 10, 2048), "must be true")
}

func TestGenerateAndVerifyProof(t *testing.T) {
	seed := []byte{0xde, 0xad, 0xbe, 0xef}

	for T := 5; T < 100; T++ {
		y_buf, proof_buf := vdf_go.GenerateVDF(seed, T, 2048)
		assert.Equal(t, true, vdf_go.VerifyVDF(seed, append(y_buf, proof_buf...), T, 2048), "failed when T = %d", T)
	}
}

func TestGenerateAndVerifyProof100(t *testing.T) {
	seed := []byte{0xde, 0xad, 0xbe, 0xef}

	for T := 101; T < 200; T++ {
		y_buf, proof_buf := vdf_go.GenerateVDF(seed, T, 2048)
		assert.Equal(t, true, vdf_go.VerifyVDF(seed, append(y_buf, proof_buf...), T, 2048), "failed when T = %d", T)
	}
}

func TestGenerateAndVerifyProof200(t *testing.T) {
	seed := []byte{0xde, 0xad, 0xbe, 0xef}

	for T := 201; T < 300; T++ {
		y_buf, proof_buf := vdf_go.GenerateVDF(seed, T, 2048)
		assert.Equal(t, true, vdf_go.VerifyVDF(seed, append(y_buf, proof_buf...), T, 2048), "failed when T = %d", T)
	}
}

func TestGenerateAndVerifyProof300(t *testing.T) {
	seed := []byte{0xde, 0xad, 0xbe, 0xef}

	for T := 301; T < 400; T++ {
		y_buf, proof_buf := vdf_go.GenerateVDF(seed, T, 2048)
		assert.Equal(t, true, vdf_go.VerifyVDF(seed, append(y_buf, proof_buf...), T, 2048), "failed when T = %d", T)
	}
}

func TestGenerateAndVerifyProof400(t *testing.T) {
	seed := []byte{0xde, 0xad, 0xbe, 0xef}

	for T := 401; T < 500; T++ {
		y_buf, proof_buf := vdf_go.GenerateVDF(seed, T, 2048)
		assert.Equal(t, true, vdf_go.VerifyVDF(seed, append(y_buf, proof_buf...), T, 2048), "failed when T = %d", T)
	}
}

func TestGenerateAndVerifyProof1000(t *testing.T) {
	seed := []byte{0xde, 0xad, 0xbe, 0xef}

	for T := 1001; T < 1010; T++ {
		y_buf, proof_buf := vdf_go.GenerateVDF(seed, T, 2048)
		assert.Equal(t, true, vdf_go.VerifyVDF(seed, append(y_buf, proof_buf...), T, 2048), "failed when T = %d", T)
	}
}

func TestRandomInput(t *testing.T) {

	for i := 0; i < 5; i++ {
		seed := make([]byte, 32)
		rand.Read(seed)

		T := 2000 + 2000*i

		y_buf, proof_buf := vdf_go.GenerateVDF(seed, T, 2048)
		assert.Equal(t, true, vdf_go.VerifyVDF(seed, append(y_buf, proof_buf...), T, 2048), "failed when T = %d", T)
	}
}
