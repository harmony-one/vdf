package main

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/harmony-one/vdf/src/vdf_go"
	"github.com/stretchr/testify/assert"
)

func TestBigNumEncoding(t *testing.T) {
	n := big.NewInt(-565721958565721958)
	s := fmt.Sprintf("%02x", vdf_go.EncodeBigIntBigEndian(n))
	assert.Equal(t, s, "fff82626e0420d489a", "they should be equal")
}

func TestEntropyFromSeed(t *testing.T) {
	seed := []byte{0xab, 0xcd}
	bytes := vdf_go.EntropyFromSeed(seed, 2048/8+2)
	s := fmt.Sprintf("%02x", bytes)
	assert.Equal(t, "aaf2bf8b0d9b3b2460493b50b1b5e55928ff4c30c9fec772c8515ced0364aeda7e570590ca8a51c2cb394edbfe7c725cbc9f9875d0bce873bc3695b95d4297d248059cf7a3e59afac18c74ea71023b26ad3ab61976e2782d2186becd79644a31030d7e6cca60e7dca796d418479bbec1167a76ca9afb933e9a66b04e70b6358a57139a45e108baf3acdef6bc021ce3a3aac077d4a0252270c3e4557acf649742fd45ad9275e5541a23f555b7c2f01dcf1a3939e9a28286e7b9b91d6f6f894f115c6d9564e80e7b49ea7736f1ae8a1c56654cbc44687a63f3218a59f8591b956b06aa7fade07145114a71a784c6bcb5de7243133a32c9a6cd6d966f37b0fabf1c6d78", s, "equal")
}

func TestDiscriminant(t *testing.T) {
	seed := []byte{0xab, 0xcd}
	n := vdf_go.CreateDiscriminant(seed, 2048)
	s := fmt.Sprintf("%02x", vdf_go.EncodeBigIntBigEndian(n))
	assert.Equal(t, "ffff550d4074f264c4db9fb6c4af4e4a1aa6d700b3cf3601388d37aea312fc9b512581a8fa6f3575ae3d34c6b12401838da34360678a2f43178c43c96a46a2bd682db7fa63085c1a65053e738b158efdc4d952c549e6891d87d2de794132869bb5cefcf28193359f182358692be7b864413ee985893565046cc165994fb18f49ca75a8ec65ba1ef7450c53210943fde31c5c553f882b5fdadd8f3c1baa85309b68bd02ba526d8a1aabe5dc0aaa483d0fe230e5c6c6165d7d79184646e2909076b0eea3926a9b17f184b61588c90e5175e3a99ab343bb97859c0cde75a607a6e46a94f95580521f8ebaeeb58e587b39434a218dbcecc5cd365932926990c84c84fe29", s, "they should be equal")
}
