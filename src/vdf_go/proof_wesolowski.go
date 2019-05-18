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
	for i := 0; i < 11; i++ {
		s1 := powers[i].a.String()
		s2 := powers[i].b.String()
		s3 := powers[i].c.String()
		fmt.Print(s1, s2, s3)
	}
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


// Get's the ith block of  2^T // B
// such that sum(get_block(i) * 2^ki) = t^T // B
func get_block(i, k, T int, B *big.Int) *big.Int {

	//(pow(2, k) * pow(2, T - k * (i + 1), B)) // B
	p1 := big.NewInt(int64(math.Pow(2, float64(k))))
	p2 := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(T - k * (i+1))), B)
	return floorDivision(new(big.Int).Mul(p1, p2), B)
}


//Optimized evalutation of h ^ (2^T // B)
func eval_optimized(identity, h *ClassGroup, B *big.Int, T, k, l int,  C map[int]*ClassGroup ) *ClassGroup{

	//k1 = k//2
	var k1 int = k / 2
	k0 := k - k1

	//x = identity
	x := CloneClassGroup(identity)

	for j := l-1; j> -1; j-- {
		//x = pow(x, pow(2, k))
		b_limit := int64(math.Pow(2, float64(k)))
		x = x.Pow(b_limit)

		//ys = {}
		ys := make([]*ClassGroup, b_limit)
		for b := int64(0); b < b_limit; b ++ {
			ys[b] = identity
		}

		//for i in range(0, math.ceil((T)/(k*l))):
		for i:= 0; i < int(math.Ceil(float64(T)/float64(k*l))) ; i++ {
			if T - k * (i*l + j + 1) < 0 {
				continue
			}

			///TODO: carefully check big.Int to int64 value conversion...might cause serious issues later
			b := get_block(i*l + j, k, T, B).Int64()

			s1 := C[i * k * l].a.String()
			s2 := C[i * k * l].b.String()
			s3 := C[i * k * l].c.String()
			fmt.Print(s1, s2, s3)

			ys[b] = ys[b].Multiply(C[i * k * l])

			s := ys[b].c.String()
			fmt.Print(s)

		}

		//for b1 in range(0, pow(2, k1)):
		for b1 := 0; b1 < int(math.Pow(float64(2), float64(k1))); b1++ {
			z := identity
			//for b0 in range(0, pow(2, k0)):
			for b0 := 0; b0 < int(math.Pow(float64(2), float64((k0)))); b0++ {
				//z *= ys[b1 * pow(2, k0) + b0]
				z = z.Multiply(ys[int64(b1) * int64(math.Pow(float64(2), float64(k0))) + int64(b0)])
			}

			//x *= pow(z, b1 * pow(2, k0))
			x = x.Multiply(z.Pow( int64(b1) * int64(math.Pow(float64(2), float64(k0)))))
		}

		//for b0 in range(0, pow(2, k0)):
		for b0 := 0; b0 < int(math.Pow(float64(2), float64(k1))); b0++ {
			z := identity
			//for b1 in range(0, pow(2, k1)):
			for b1 := 0; b1 < int(math.Pow(float64(2), float64(k1))); b1++ {
				//z *= ys[b1 * pow(2, k0) + b0]
				z = z.Multiply(ys[int64(b1) * int64(math.Pow(float64(2), float64(k0))) + int64(b0)])
			}
			//x *= pow(z, b0)
			x = x.Multiply(z.Pow(int64(b0)))
		}
	}


	return x
}


//generate y = x ^ (2 ^T) and pi
func generate_proof(identity, x, y *ClassGroup, T, k, l int, powers map[int]*ClassGroup) *ClassGroup {
	//x_s = x.serialize()
	x_s := x.Serialize()

	//y_s = y.serialize()
	y_s := y.Serialize()

	B := hash_prime(x_s, y_s)

	//s := fmt.Sprintf("%02x", EncodeBigIntBigEndian(B))
	//fmt.Print(s)

	proof := eval_optimized(identity, x, B, T, k, l, powers)

	s1 := proof.a.String()
	s2 := proof.b.String()
	s3 := proof.c.String()
	fmt.Print(s1, s2, s3)

	ya := y.a.String()
	yb := y.b.String()
	yc := y.c.String()
	fmt.Print(ya,yb,yc)

	return proof
}