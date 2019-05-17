package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"vdf_go"
)

func TestDiscriminant(t *testing.T) {

	//L, k, _ := approximateParameters(20000000)

	seed := []byte{0xab,0xcd}
	/*
		bytes := entropy_from_seed(seed, 258)
		for i := 0; i< len(bytes); i++ {
			fmt.Printf("%x,", bytes[i])
		}
		fmt.Print("\n")
	*/
	n := vdf_go.CreateDiscriminant(seed, 2048)
	s := fmt.Sprintf("%02x", vdf_go.EncodeBigIntBigEndian(n))
	assert.Equal(t, s, "ffde47c49afffffffd1c", "they should be equal")

}

