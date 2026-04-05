package vdf_go

import (
	"math/big"
)

func absBigForGCD(x *big.Int) *big.Int {
	if x.Sign() >= 0 {
		return x
	}
	return new(big.Int).Neg(x)
}

// wrapper around big.Int GCD to allow all input values for GCD
// as Golang big.Int GCD requires both a, b > 0
// If a == b == 0, GCD sets r = 0.
// If a == 0 and b != 0, GCD sets r = |b|
// If a != 0 and b == 0, GCD sets r = |a|
// Otherwise r = GCD(|a|, |b|)
func allInputValueGCD(a, b *big.Int) (r *big.Int) {
	if a.Sign() == 0 {
		return new(big.Int).Abs(b)
	}

	if b.Sign() == 0 {
		return new(big.Int).Abs(a)
	}

	return new(big.Int).GCD(nil, nil, absBigForGCD(a), absBigForGCD(b))
}

// Solve ax == b mod m for x.
// Return s, t where x = s + k * t for integer k yields all solutions.
func SolveMod(a, b, m *big.Int) (s, t *big.Int, solvable bool) {
	absA := absBigForGCD(a)
	absM := absBigForGCD(m)
	d := new(big.Int)
	g := new(big.Int).GCD(d, nil, absA, absM)
	if a.Sign() < 0 {
		d.Neg(d)
	}

	var r big.Int
	q, rem := new(big.Int).QuoRem(b, g, &r)

	//TODO: replace with utils.GetLogInstance().Error(...)
	//if r != 0:
	if rem.Sign() != 0 {
		//panic(fmt.Sprintf("no solution to %s x = %s mod %s", a.String(), b.String(), m.String()))
		return nil, nil, false
	}

	//assert b == q * g
	//return (q * d) % m, m // g
	q.Mul(q, d)
	s = q.Mod(q, absM)
	t = new(big.Int).Quo(absM, g)
	return s, t, true
}
