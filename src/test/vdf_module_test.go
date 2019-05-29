package main

import (
	"crypto/rand"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
	"vdf_go"
)

func TestVDFModule(t *testing.T) {

	input := [32]byte{}
	rand.Read(input[:])

	vdf := vdf_go.New(10000, input)

	outputChannel := vdf.GetOutputChannel()

	start := time.Now()

	vdf.Execute()

	duration := time.Now().Sub(start)

	log.Println(fmt.Sprintf("VDF computation finished, time spent %s", duration.String()))

	output := <-outputChannel

	assert.Equal(t, true, vdf.Verify(output), "failed verifying proof")

}
