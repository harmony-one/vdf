package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"math/big"
	"regexp"
	"runtime"
	"testing"
	"time"
	"vdf_go"
)

func RepeatedSquare(x *vdf_go.ClassGroup, k int) *vdf_go.ClassGroup {
	defer timeTrack(time.Now())

	for i := 0; i < k; i++ {
		x = x.Square()
	}

	return x
}

func RepeatedSquareSlow(x *vdf_go.ClassGroup, k int) *vdf_go.ClassGroup {
	defer timeTrack(time.Now())

	for i := 0; i < k; i++ {
		x = x.SquareUsingMultiply()
	}

	return x
}

func TestTwoSquarePerformance(t *testing.T) {
	for k := 0; k < 10; k++ {
		seed := make([]byte, 32)
		rand.Read(seed)
		D := vdf_go.CreateDiscriminant(seed, 2048)
		x := vdf_go.NewClassGroupFromAbDiscriminant(big.NewInt(2), big.NewInt(1), D)

		y := vdf_go.CloneClassGroup(x)
		y1 := vdf_go.CloneClassGroup(x)

		y = RepeatedSquare(y, 5000)
		y1 = RepeatedSquareSlow(y1, 5000)

		assert.Equal(t, true, y.Equal(y1), "k=%d, seed=%s", k, hex.EncodeToString(seed))
		log.Print(fmt.Sprintf("Test case %d good", k))
	}
}

func timeTrack(start time.Time) {
	elapsed := time.Since(start)

	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a function object this functions parent.
	funcObj := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path).
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	log.Println(fmt.Sprintf("%s took %s", name, elapsed))
}
