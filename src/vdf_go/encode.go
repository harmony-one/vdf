package vdf_go

import "math/big"

//inplace encoding using two's complementary
func two_s_complement_encoding(buf []byte, bytes_size int) []byte {
	//two's complement carry
	var carry uint8 = 1

	//use one additional byte for signing
	for i := len(buf) - 1; i >= len(buf)-bytes_size; i-- {
		thisdigit := uint8(buf[i])
		thisdigit = thisdigit ^ 0xff

		if thisdigit == 0xff {
			if carry == 1 {
				thisdigit = 0
				carry = 1
			} else {
				carry = 0
			}
		} else {
			thisdigit = thisdigit + carry
			carry = 0
		}

		buf[i] = thisdigit
	}

	//put all remaining leading bytes to 0
	for i := len(buf) - bytes_size - 1; i >= 0; i-- {
		buf[i] = 0xff
	}

	return buf
}

func EncodeBigIntBigEndian(a *big.Int) []byte {
	int_size_bits := a.BitLen()
	int_size := (int_size_bits + 16) >> 3

	buf := make([]byte, int_size)
	a_bytes := a.Bytes()
	copy(buf[int_size-len(a_bytes):], a_bytes)

	//encode the negative number
	if a.Sign() == -1 {
		two_s_complement_encoding(buf, len(a_bytes))
	}

	return buf
}
