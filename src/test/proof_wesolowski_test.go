package main

import (
	"math/big"
	"testing"
	"vdf_go"
)

func Test1(t *testing.T) {
	seed := []byte{0xde, 0xad, 0xbe, 0xef}
	D := vdf_go.CreateDiscriminant(seed, 2048)

	X := vdf_go.NewClassGroupFromAbDiscriminant(big.NewInt(2), big.NewInt(1), D)
	vdf_go.Create_proof_of_time_wesolowski(D, X, 10, 2048)

}
