package vdf_go

import (
	"math/big"
)

//Return r, s, t such that gcd(a, b) = r = a * s + b * t
func extendedGCD(a, b *big.Int) (r, s, t *big.Int) {
	//r0, r1 = a, b
	r0 := new(big.Int).Set(a)
	r1 := new(big.Int).Set(b)

	//s0, s1, t0, t1 = 1, 0, 0, 1
	s0 := big.NewInt(1)
	s1 := big.NewInt(0)
	t0 := big.NewInt(0)
	t1 := big.NewInt(1)

	//if r0 > r1:
	//r0, r1, s0, s1, t0, t1 = r1, r0, t0, t1, s0, s1
	if r0.Cmp(r1) == 1 {
		oldR0 := new(big.Int).Set(r0)
		r0 = r1
		r1 = oldR0
		oldS0 := new(big.Int).Set(s0)
		s0 = t0
		oldS1 := new(big.Int).Set(s1)
		s1 = t1
		t0 = oldS0
		t1 = oldS1
	}

	//while r1 > 0:
	for r1.Sign() == 1 {
		//q, r = divmod(r0, r1)
		r := big.NewInt(1)
		bb := new(big.Int).Set(b)
		q, r := bb.DivMod(r0, r1, r)

		//r0, r1, s0, s1, t0, t1 = r1, r, s1, s0 - q * s1, t1, t0 - q * t1
		r0 = r1
		r1 = r
		oldS0 := new(big.Int).Set(s0)
		s0 = s1
		s1 = new(big.Int).Sub(oldS0, new(big.Int).Mul(q, s1))
		oldT0 := new(big.Int).Set(t0)
		t0 = t1
		t1 = new(big.Int).Sub(oldT0, new(big.Int).Mul(q, t1))

	}
	return r0, s0, t0
}

//wrapper around big.Int GCD to allow all input values for GCD
//as Golang big.Int GCD requires both a, b > 0
//If a == b == 0, GCD sets r = 0.
//If a == 0 and b != 0, GCD sets r = |b|
//If a != 0 and b == 0, GCD sets r = |a|
//Otherwise r = GCD(|a|, |b|)
func allInputValueGCD(a, b *big.Int) (r *big.Int) {
	if a.Sign() == 0 {
		return new(big.Int).Abs(b)
	}

	if b.Sign() == 0 {
		return new(big.Int).Abs(a)
	}

	return new(big.Int).GCD(nil, nil, new(big.Int).Abs(a), new(big.Int).Abs(b))
}

//Solve ax == b mod m for x.
//Return s, t where x = s + k * t for integer k yields all solutions.
func SolveMod(a, b, m *big.Int) (s, t *big.Int, solvable bool) {
	//g, d, e = extended_gcd(a, m)
	//TODO: golang 1.x big.int GCD requires both a > 0 and m > 0, so we can't use it :(
	//d := big.NewInt(0)
	//e := big.NewInt(0)
	//g := new(big.Int).GCD(d, e, a, m)
	g, d, _ := extendedGCD(a, m)

	//q, r = divmod(b, g)
	r := big.NewInt(1)
	bb := new(big.Int).Set(b)
	q, r := bb.DivMod(b, g, r)

	//TODO: replace with utils.GetLogInstance().Error(...)
	//if r != 0:
	if r.Cmp(big.NewInt(0)) != 0 {
		//panic(fmt.Sprintf("no solution to %s x = %s mod %s", a.String(), b.String(), m.String()))
		return nil, nil, false
	}

	//assert b == q * g
	//return (q * d) % m, m // g
	q.Mul(q, d)
	s = q.Mod(q, m)
	t = floorDivision(m, g)
	return s, t, true
}
