package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"math/big"
	"testing"
	"vdf_go"
)

func TestSquare3(t *testing.T) {
	//x := vdf_go.NewClassGroup(big.NewInt(565721958), big.NewInt(-740), big.NewInt(4486780496))
	//x := vdf_go.NewClassGroup(big.NewInt(93), big.NewInt(109), big.NewInt(32))
	//x := vdf_go.NewClassGroup(big.NewInt(195751), big.NewInt(1212121), big.NewInt(1876411))

	for k :=0 ; k< 1000; k++ {
		seed := make([]byte, 32)
		rand.Read(seed)
		D := vdf_go.CreateDiscriminant(seed, 2048)
		x := vdf_go.NewClassGroupFromAbDiscriminant(big.NewInt(2), big.NewInt(1), D)

		y  := vdf_go.CloneClassGroup(x)
		y1 := vdf_go.CloneClassGroup(x)

		for i := 0; i< 1000 + k; i++ {
			y  = y.Square()
			y1 = y1.Square1()
		}

		assert.Equal(t, true, y.Equal(y1), "k=%d, seed=%s",k, hex.EncodeToString(seed))
		log.Print(fmt.Sprintf("Test case %d good", k))
	}
}
