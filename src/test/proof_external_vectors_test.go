package main

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/harmony-one/vdf/src/vdf_go"
	"github.com/stretchr/testify/assert"
)

// Source vectors (Wesolowski):
// https://github.com/poanetwork/vdf/blob/master/wesolowski.csv
// (fetched on 2026-04-05, first three rows).
func TestRustReferenceWesolowskiCanaryVectors(t *testing.T) {
	type vector struct {
		seedHex         string
		iterations      int
		outputSHA256Hex string
	}

	vectors := []vector{
		{
			seedHex:         "acc390feb1fbbe70d6a7ad2203a8b8c3c93e52d9782886124606686bc5716fed",
			iterations:      66,
			outputSHA256Hex: "3d72a0852b85b87e466c1268233ae02c31295497dd88b92b63561141b1b68b29",
		},
		{
			seedHex:         "c86348990a772d094da4b5d4a3014a2d981c4b36fa1df672456d98d69bfb3f9e",
			iterations:      68,
			outputSHA256Hex: "4ea8c6df1a0a85e12f82320f631f4f276a5a373a0de4f65d7193ae1f1e66d132",
		},
		{
			seedHex:         "7bd5c3dba392a0ef4b02052200748d5065e4d5d2f3bbf11ca5edeaccfc44f7aa",
			iterations:      70,
			outputSHA256Hex: "ccf96d2bd863cf7fa5ab992be2d6dd2f5f500da069ac07b34c06f6c367e67e75",
		},
	}

	for _, tc := range vectors {
		tc := tc
		t.Run(tc.seedHex, func(t *testing.T) {
			seed, err := hex.DecodeString(tc.seedHex)
			assert.NoError(t, err)

			yBuf, proofBuf := vdf_go.GenerateVDF(seed, tc.iterations, 2048)
			assert.NotNil(t, yBuf)
			assert.NotNil(t, proofBuf)

			output := append(yBuf, proofBuf...)

			assert.True(t, vdf_go.VerifyVDF(seed, output, tc.iterations, 2048))

			digest := sha256.Sum256(output)
			assert.Equal(t, tc.outputSHA256Hex, hex.EncodeToString(digest[:]))
		})
	}
}
