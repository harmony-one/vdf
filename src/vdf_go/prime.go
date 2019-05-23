package vdf_go

import (
	"log"
	"math"
	"math/big"
)

func primeLessThanN(num int) []int {
	//initialized to false
	sieve := make([]bool, num+1)

	for i := 3; i <= int(math.Floor(math.Sqrt(float64(num)))); i += 2 {
		if sieve[i] == false {
			for j := i * 2; j <= num; j += i {
				sieve[j] = true // cross
			}
		}
	}

	primes := make([]int, 0, num)
	for i := 3; i <= num; i += 2 {
		if sieve[i] == false {
			primes = append(primes, i)
		}
	}

	return primes
}

//testing
func checkArrayEqual(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func testIntLessThan(num int) {
	var refPrimes = make([]int, 0, num)
	for i := 3; i < num; i += 2 {
		if big.NewInt(int64(i)).ProbablyPrime(1) {
			refPrimes = append(refPrimes, i)
		}
	}

	var primes = primeLessThanN(num)
	log.Printf("%v ", primes)

	if checkArrayEqual(refPrimes, primes) {
		log.Printf("OK")
	} else {
		log.Printf("ERROR")
	}
}

func main() {
	testIntLessThan(1 << 16)
}
