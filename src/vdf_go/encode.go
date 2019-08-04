package vdf_go

import "math/big"

var bigOne = big.NewInt(1)

func decodeTwosComplement(bytes []byte) *big.Int {
	if bytes[0]&0x80 == 0 {
		// non-negative
		return new(big.Int).SetBytes(bytes)
	}
	setyb := make([]byte, len(bytes))
	for i := range bytes {
		setyb[i] = bytes[i] ^ 0xff
	}
	n := new(big.Int).SetBytes(setyb)
	return n.Sub(n.Neg(n), bigOne)
}

func encodeTwosComplement(n *big.Int) []byte {
	if n.Sign() > 0 {
		bytes := n.Bytes()
		if bytes[0]&0x80 == 0 {
			return bytes
		}
		// add one more byte for positive sign
		buf := make([]byte, len(bytes)+1)
		copy(buf[1:], bytes)
		return buf
	}
	if n.Sign() < 0 {
		// A negative number has to be converted to two's-complement form. So we
		// invert and subtract 1. If the most-significant-bit isn't set then
		// we'll need to pad the beginning with 0xff in order to keep the number
		// negative.
		nMinus1 := new(big.Int).Neg(n)
		nMinus1.Sub(nMinus1, bigOne)
		bytes := nMinus1.Bytes()
		if len(bytes) == 0 {
			// sneaky -1 value
			return []byte{0xff}
		}
		for i := range bytes {
			bytes[i] ^= 0xff
		}
		if bytes[0]&0x80 != 0 {
			return bytes
		}
		// add one more byte for negative sign
		buf := make([]byte, len(bytes)+1)
		buf[0] = 0xff
		copy(buf[1:], bytes)
		return buf
	}
	return []byte{}
}

func signBitFill(bytes []byte, targetLen int) []byte {
	if len(bytes) >= targetLen {
		return bytes
	}
	buf := make([]byte, targetLen)
	offset := targetLen - len(bytes)
	if bytes[0]&0x80 != 0 {
		for i := 0; i < offset; i++ {
			buf[i] = 0xff
		}
	}
	copy(buf[offset:], bytes)
	return buf
}

func EncodeBigIntBigEndian(a *big.Int) []byte {
	int_size_bits := a.BitLen()
	int_size := (int_size_bits + 16) >> 3

	return signBitFill(encodeTwosComplement(a), int_size)
}
