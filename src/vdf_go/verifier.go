package vdf_go

import "math/big"

type verifierState struct {
	discriminant   *big.Int
	base           *ClassGroup
	baseSerialized []byte
	intSize        int
}

func serializeReducedClassGroup(group *ClassGroup, intSize int) []byte {
	buf := make([]byte, intSize*2)
	fillTwosComplement(buf[:intSize], group.a)
	fillTwosComplement(buf[intSize:], group.b)
	return buf
}

func newVerifierState(seed []byte, intSizeBits int) *verifierState {
	intSize := (intSizeBits + 16) >> 4
	discriminant := CreateDiscriminant(seed, intSizeBits)
	base := NewClassGroupFromAbDiscriminant(big.NewInt(2), big.NewInt(1), discriminant)

	return &verifierState{
		discriminant:   discriminant,
		base:           base,
		baseSerialized: serializeReducedClassGroup(base, intSize),
		intSize:        intSize,
	}
}

func multiExpReduced(a *ClassGroup, eA *big.Int, b *ClassGroup, eB *big.Int) *ClassGroup {
	maxBits := eA.BitLen()
	if eB.BitLen() > maxBits {
		maxBits = eB.BitLen()
	}
	if maxBits == 0 {
		return a.identity()
	}

	id := a.identity()
	powA := [4]*ClassGroup{id, a}
	powB := [4]*ClassGroup{id, b}

	powA[2] = a.Square()
	if powA[2] == nil {
		return nil
	}
	powA[3] = powA[2].multiplyReduced(a)
	if powA[3] == nil {
		return nil
	}

	powB[2] = b.Square()
	if powB[2] == nil {
		return nil
	}
	powB[3] = powB[2].multiplyReduced(b)
	if powB[3] == nil {
		return nil
	}

	table := [16]*ClassGroup{}
	table[1] = powB[1]
	table[2] = powB[2]
	table[3] = powB[3]
	table[4] = powA[1]
	table[8] = powA[2]
	table[12] = powA[3]

	for ea := 1; ea < 4; ea++ {
		for eb := 1; eb < 4; eb++ {
			idx := (ea << 2) | eb
			table[idx] = powA[ea].multiplyReduced(powB[eb])
			if table[idx] == nil {
				return nil
			}
		}
	}

	windows := (maxBits + 1) >> 1
	var acc *ClassGroup
	started := false

	for w := windows - 1; w >= 0; w-- {
		if started {
			acc = acc.Square()
			if acc == nil {
				return nil
			}
			acc = acc.Square()
			if acc == nil {
				return nil
			}
		}

		bit := w << 1
		eaDigit := (eA.Bit(bit+1) << 1) | eA.Bit(bit)
		ebDigit := (eB.Bit(bit+1) << 1) | eB.Bit(bit)
		key := (eaDigit << 2) | ebDigit
		if key == 0 {
			continue
		}

		factor := table[key]
		if !started {
			acc = factor
			started = true
			continue
		}

		acc = acc.multiplyReduced(factor)
		if acc == nil {
			return nil
		}
	}

	if !started {
		return id
	}

	return acc
}

func verifyProofBlob(state *verifierState, proofBlob []byte, iterations int) bool {
	ySerialized := proofBlob[:2*state.intSize]
	y, _ := NewClassGroupFromBytesDiscriminant(ySerialized, state.discriminant)
	proof, _ := NewClassGroupFromBytesDiscriminant(proofBlob[2*state.intSize:], state.discriminant)

	B := hashPrime(state.baseSerialized, ySerialized)
	r := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(iterations)), B)

	z := multiExpReduced(proof, B, state.base, r)
	return z != nil && z.equalReduced(y)
}
