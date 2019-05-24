package main

import (
	"bufio"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"os"
	"strconv"
	"testing"
	"vdf_go"
)

func TestCreateProofCSV(t *testing.T) {

	csvFile, _ := os.Open("wesolowski.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	for i := 0; ; i++ {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		seed, _ := hex.DecodeString(line[0])
		T, _ := strconv.Atoi(line[1])
		P := line[2]

		y_buf, proof_buf := vdf_go.GenerateVDF(seed, T, 2048)
		proof := hex.EncodeToString(append(y_buf, proof_buf...))
		assert.Equal(t, P, proof, "iteration %d", T)
		log.Print(fmt.Sprintf("Test case %d good, iteration = %d", i, T))

	}
}

func TestVerifyProofCSV(t *testing.T) {

	csvFile, _ := os.Open("wesolowski.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	for i := 0; ; i++ {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		seed, _ := hex.DecodeString(line[0])
		T, _ := strconv.Atoi(line[1])
		P, _ := hex.DecodeString(line[2])

		result := vdf_go.VerifyVDF(seed, P, T, 2048)
		assert.Equal(t, true, result, "iteration %d", T)
		log.Print(fmt.Sprintf("Test case %d good, iteration = %d", i, T))

	}
}
