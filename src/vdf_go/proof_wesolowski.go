package vdf_go

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

//Creates L and k parameters from papers, based on how many iterations need to be
//performed, and how much memory should be used.
func approximateParameters(T int) (int, int, int) {
	//log_memory = math.log(10000000, 2)
	log_memory := math.Log(10000000) / math.Log(2)
	log_T := math.Log(float64(T)) / math.Log(2)
	L := 1

	if (log_T - log_memory > 0){
		L = int(math.Ceil(math.Pow(2, log_memory - 20)))
	}

	// Total time for proof: T/k + L * 2^(k+1)
	// To optimize, set left equal to right, and solve for k
	// k = W(T * log(2) / (2 * L))  / log(2), where W is the product log function
	// W can be approximated by log(x) - log(log(x)) + 0.25
	intermediate := float64(T) * math.Log(2) / float64(2 * L)
	k := int(math.Max(math.Round(math.Log(intermediate) - math.Log(math.Log(intermediate)) + 0.25), 1))

	// 1/w is the approximate proportion of time spent on the proof
	w := int(math.Floor(float64(T) / (float64(T)/float64(k) + float64(L) * math.Pow(2, float64(k + 1)))) - 2)

	return L, k,  w
}

func iterateSquarings(x *ClassGroup, powers_to_calculate []int) map[int]*ClassGroup {
	powers_calculated := make(map[int]*ClassGroup)

	previous_power := 0
	currX := CloneClassGroup(x)
	for _, current_power := range powers_to_calculate {

		for i:=0; i < current_power - previous_power; i++ {
			currX = currX.Pow(2)
			s := fmt.Sprintf("%02x", EncodeBigIntAbs(currX.c, 256))
			fmt.Print(s)
		}

		previous_power = current_power
		powers_calculated[current_power] = currX

		//s := fmt.Sprintf("%02x", EncodeBigIntBigEndian(currX.c))
		//fmt.Print(s)

		//currX = CloneClassGroup(currX)
	}


	return powers_calculated
}

func Create_proof_of_time_wesolowski(discriminant *big.Int, x *ClassGroup, iterations, int_size_bits int) map[int]*ClassGroup {
	L, k, _ := approximateParameters(iterations)

	loopCount := int(math.Ceil(float64(iterations)/float64(k*L)))
	powers_to_calculate := make([]int, loopCount + 2)

	for i := 0; i < loopCount + 1; i++ {
		powers_to_calculate[i] = i * k * L
	}

	powers_to_calculate[loopCount + 1] = loopCount

	powers := iterateSquarings(x, powers_to_calculate)

	y := powers[loopCount]

	/*
	a := fmt.Sprintf("%02x", EncodeBigIntBigEndian(y.a))
	b := fmt.Sprintf("%02x", EncodeBigIntBigEndian(y.b))
	c := fmt.Sprintf("%02x", EncodeBigIntBigEndian(y.c))
	fmt.Print(a, b, c)
	*/

	identity := IdentityForDiscriminant(discriminant)

	generate_proof(identity, x, y, iterations, k, L, powers)

	return powers

}

// Creates a random prime based on input x, y
func hash_prime(x, y []byte)  *big.Int {

	var j uint64 = 0

	jBuf := make([]byte, 8)
	z := new(big.Int)
	for {
		binary.BigEndian.PutUint64(jBuf, j)
		s := append([]byte("prime"), jBuf...)
		s  = append(s, x... )
		s  = append(s, y... )

		checkSum := sha256.Sum256(s[:])
		z.SetBytes(checkSum[:16])
		aa := fmt.Sprintf("%02x", EncodeBigIntBigEndian(z))
		fmt.Print(aa)
		if z.ProbablyPrime(1) {
			return z
		}
		j++

	}
}

//Optimized evalutation of h ^ (2^T // B)
func eval_optimized(identity, h *ClassGroup, B *big.Int, T, k, l int,  C map[int]*ClassGroup ) *big.Int{
	//k1 = k//2
	k1 := floorDivision(k, big.NewInt(2))

	//k0 = k - k1
	k0 : = new(big.Int).Sub(k, k1)

	//x = identity
	x := CloneClassGroup(identity)

	//ys = {}
	p2K = int(math.Pow(2, float64(k)))
	ys := make([]*ClassGroup, p2K)
	for j := l-1; j> -1; j-- {
		//x = pow(x, pow(2, k))
		p2K = int(math.Pow(2, float64(k)))
		x = x.Pow(p2K)

		//b_limit = pow(2, k)
		b_limit := p2K

		for b := 0; b < b_limit; b ++ {
			ys[b] = identity
		}
	}


}


//generate y = x ^ (2 ^T) and pi
func generate_proof(identity, x, y *ClassGroup, T, k, l int, powers map[int]*ClassGroup) *big.Int {
	//x_s = x.serialize()
	x_s := x.Serialize()

	//y_s = y.serialize()
	y_s := y.Serialize()


	B := hash_prime(x_s, y_s)

	#s := fmt.Sprintf("%02x", EncodeBigIntBigEndian(B))
	#fmt.Print(s)


	return B
}