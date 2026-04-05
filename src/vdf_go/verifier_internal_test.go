package vdf_go

import (
	"math/big"
	"testing"
)

func TestMultiExpReducedMatchesSeparateExponentiation(t *testing.T) {
	seed := []byte{
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
	}

	discriminant := CreateDiscriminant(seed, 2048)
	base := NewClassGroupFromAbDiscriminant(big.NewInt(2), big.NewInt(1), discriminant)
	other := base.Pow(17)
	if other == nil {
		t.Fatal("failed to build comparison class group")
	}

	cases := []struct {
		a string
		b string
	}{
		{"0", "0"},
		{"1", "1"},
		{"123456789abcdef", "fedcba987654321"},
		{"ffffffffffffffffffffffffffffffff", "123456789abcdef123456789abcdef"},
	}

	for _, tc := range cases {
		eA, ok := new(big.Int).SetString(tc.a, 16)
		if !ok {
			t.Fatalf("invalid exponent A: %s", tc.a)
		}
		eB, ok := new(big.Int).SetString(tc.b, 16)
		if !ok {
			t.Fatalf("invalid exponent B: %s", tc.b)
		}

		got := multiExpReduced(other, eA, base, eB)
		if got == nil {
			t.Fatalf("multiExpReduced returned nil for %s/%s", tc.a, tc.b)
		}

		left := other.bigPowReduced(eA)
		right := base.bigPowReduced(eB)
		if left == nil || right == nil {
			t.Fatalf("separate exponentiation returned nil for %s/%s", tc.a, tc.b)
		}

		want := left.multiplyReduced(right)
		if want == nil {
			t.Fatalf("separate multiplication returned nil for %s/%s", tc.a, tc.b)
		}

		if !got.equalReduced(want) {
			t.Fatalf("mismatch for exponents %s/%s", tc.a, tc.b)
		}
	}
}
