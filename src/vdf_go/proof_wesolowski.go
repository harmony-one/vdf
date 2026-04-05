package vdf_go

import (
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math"
	"math/big"
	"regexp"
	"runtime"
	"sync/atomic"
	"time"
)

var (
	timingEnabled     atomic.Bool
	runtimeFuncNameRE = regexp.MustCompile(`^.*\.(.*)$`)
	primePrefix       = []byte("prime")
)

// Creates L and k parameters from papers, based on how many iterations need to be
// performed, and how much memory should be used.
func approximateParameters(T int) (int, int, int) {
	//log_memory = math.log(10000000, 2)
	log_memory := math.Log(10000000) / math.Log(2)
	log_T := math.Log(float64(T)) / math.Log(2)
	L := 1

	if log_T-log_memory > 0 {
		L = int(math.Ceil(math.Pow(2, log_memory-20)))
	}

	// Total time for proof: T/k + L * 2^(k+1)
	// To optimize, set left equal to right, and solve for k
	// k = W(T * log(2) / (2 * L))  / log(2), where W is the product log function
	// W can be approximated by log(x) - log(log(x)) + 0.25
	intermediate := float64(T) * math.Log(2) / float64(2*L)
	k := int(math.Max(math.Round(math.Log(intermediate)-math.Log(math.Log(intermediate))+0.25), 1))

	// 1/w is the approximate proportion of time spent on the proof
	w := int(math.Floor(float64(T)/(float64(T)/float64(k)+float64(L)*math.Pow(2, float64(k+1)))) - 2)

	return L, k, w
}

func repeatedSquare(x *ClassGroup, count int, stop <-chan struct{}) *ClassGroup {
	if stop == nil {
		for i := 0; i < count; i++ {
			x = x.Square()
			if x == nil {
				return nil
			}
		}

		return x
	}

	for i := 0; i < count; i++ {
		select {
		case <-stop:
			return nil
		default:
		}

		x = x.Square()
		if x == nil {
			return nil
		}
	}

	return x
}

func iterateSquarings(x *ClassGroup, step, iterations int, stop <-chan struct{}) ([]*ClassGroup, *ClassGroup) {
	if iterations <= 0 {
		return []*ClassGroup{x}, x
	}

	loopCount := int(math.Ceil(float64(iterations) / float64(step)))
	powersCalculated := make([]*ClassGroup, loopCount)
	powersCalculated[0] = x

	currX := x
	for i := 1; i < loopCount; i++ {
		currX = repeatedSquare(currX, step, stop)
		if currX == nil {
			return nil, nil
		}
		powersCalculated[i] = currX
	}

	remaining := iterations - (loopCount-1)*step
	y := repeatedSquare(currX, remaining, stop)
	if y == nil {
		return nil, nil
	}

	return powersCalculated, y
}

func GenerateVDF(seed []byte, iterations, int_size_bits int) ([]byte, []byte) {
	return GenerateVDFWithStopChan(seed, iterations, int_size_bits, nil)
}

func GenerateVDFWithStopChan(seed []byte, iterations, int_size_bits int, stop <-chan struct{}) ([]byte, []byte) {
	if timingEnabled.Load() {
		defer timeTrack(time.Now())
	}

	D := CreateDiscriminant(seed, int_size_bits)
	x := NewClassGroupFromAbDiscriminant(big.NewInt(2), big.NewInt(1), D)

	y, proof := calculateVDF(D, x, iterations, int_size_bits, stop)

	if (y == nil) || (proof == nil) {
		return nil, nil
	} else {
		return y.Serialize(), proof.Serialize()
	}
}

func VerifyVDF(seed, proof_blob []byte, iterations, int_size_bits int) bool {
	if timingEnabled.Load() {
		defer timeTrack(time.Now())
	}

	return verifyProofBlob(newVerifierState(seed, int_size_bits), proof_blob, iterations)
}

// Creates a random prime based on input x, y
func hashPrime(x, y []byte) *big.Int {
	var j uint64
	buf := make([]byte, len(primePrefix)+8+len(x)+len(y))
	copy(buf, primePrefix)
	copy(buf[len(primePrefix)+8:], x)
	copy(buf[len(primePrefix)+8+len(x):], y)
	z := new(big.Int)
	for {
		binary.BigEndian.PutUint64(buf[len(primePrefix):], j)
		checkSum := sha256.Sum256(buf)
		z.SetBytes(checkSum[:16])

		if z.ProbablyPrime(1) {
			return z
		}
		j++

	}
}

// Get's the ith block of  2^T // B
// such that sum(get_block(i) * 2^ki) = t^T // B
func getBlock(i, k, T int, B, twoPowK *big.Int) *big.Int {
	p2 := new(big.Int).Exp(bigTwoValue, big.NewInt(int64(T-k*(i+1))), B)
	p2.Mul(p2, twoPowK)
	return p2.Quo(p2, B)
}

// Optimized evalutation of h ^ (2^T // B)
func evalOptimized(identity, h *ClassGroup, B *big.Int, T, k, l int, C []*ClassGroup) *ClassGroup {
	//k1 = k//2
	var k1 int = k / 2
	k0 := k - k1
	outerIterations := int(math.Ceil(float64(T) / float64(k*l)))
	bLimit := 1 << uint(k)
	b1Limit := 1 << uint(k1)
	b0Limit := 1 << uint(k0)
	b0Stride := int64(b0Limit)
	twoPowK := new(big.Int).Lsh(new(big.Int).Set(bigOne), uint(k))

	//x = identity
	x := identity
	ys := make([]*ClassGroup, bLimit)

	for j := l - 1; j > -1; j-- {
		x = repeatedSquare(x, k, nil)
		if x == nil {
			return nil
		}

		for b := range ys {
			ys[b] = identity
		}

		for i := 0; i < outerIterations; i++ {
			if T-k*(i*l+j+1) < 0 {
				continue
			}

			b := int(getBlock(i*l+j, k, T, B, twoPowK).Int64())
			ys[b] = ys[b].Multiply(C[i])
			if ys[b] == nil {
				return nil
			}
		}

		for b1 := 0; b1 < b1Limit; b1++ {
			z := identity
			for b0 := 0; b0 < b0Limit; b0++ {
				z = z.Multiply(ys[b1*b0Limit+b0])
				if z == nil {
					return nil
				}
			}

			c := z.Pow(int64(b1) * b0Stride)
			if c == nil {
				return nil
			}
			x = x.Multiply(c)
			if x == nil {
				return nil
			}
		}

		for b0 := 0; b0 < b0Limit; b0++ {
			z := identity
			for b1 := 0; b1 < b1Limit; b1++ {
				z = z.Multiply(ys[b1*b0Limit+b0])
				if z == nil {
					return nil
				}
			}
			d := z.Pow(int64(b0))
			if d == nil {
				return nil
			}
			x = x.Multiply(d)
			if x == nil {
				return nil
			}
		}
	}

	return x
}

// generate y = x ^ (2 ^T) and pi
func generateProof(identity, x, y *ClassGroup, T, k, l int, powers []*ClassGroup) *ClassGroup {
	//x_s = x.serialize()
	x_s := x.Serialize()

	//y_s = y.serialize()
	y_s := y.Serialize()

	B := hashPrime(x_s, y_s)

	proof := evalOptimized(identity, x, B, T, k, l, powers)

	return proof
}

func calculateVDF(discriminant *big.Int, x *ClassGroup, iterations, int_size_bits int, stop <-chan struct{}) (y, proof *ClassGroup) {
	L, k, _ := approximateParameters(iterations)

	powers, y := iterateSquarings(x, k*L, iterations, stop)
	if powers == nil || y == nil {
		return nil, nil
	}

	identity := IdentityForDiscriminant(discriminant)

	proof = generateProof(identity, x, y, iterations, k, L, powers)

	return y, proof
}

func verifyProof(x, y, proof *ClassGroup, T int) bool {
	//x_s = x.serialize()
	x_s := x.Serialize()

	//y_s = y.serialize()
	y_s := y.Serialize()

	B := hashPrime(x_s, y_s)

	r := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(T)), B)

	piB := proof.BigPow(B)
	if piB == nil {
		return false
	}

	xR := x.BigPow(r)
	if xR == nil {
		return false
	}

	z := piB.Multiply(xR)
	if (z != nil) && (z.Equal(y)) {
		return true
	} else {
		return false
	}
}

func SetTimingEnabled(enabled bool) {
	timingEnabled.Store(enabled)
}

func timeTrack(start time.Time) {
	elapsed := time.Since(start)

	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a function object this functions parent.
	funcObj := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path).
	name := runtimeFuncNameRE.ReplaceAllString(funcObj.Name(), "$1")

	log.Printf("%s took %s", name, elapsed)
}
