package vdf_go

import (
	"math/big"
)

type ClassGroup struct {
	a *big.Int
	b *big.Int
	c *big.Int
	d *big.Int
}

func CloneClassGroup(cg *ClassGroup) *ClassGroup {
	return &ClassGroup{a: cg.a, b: cg.b, c: cg.c}
}

func NewClassGroup(a, b, c *big.Int) *ClassGroup {
	return &ClassGroup{a: a, b: b, c: c}
}

func NewClassGroupFromAbDiscriminant(a, b, discriminant *big.Int) *ClassGroup {
	//z = b*b-discriminant
	z := new(big.Int).Sub(new(big.Int).Mul(b, b), discriminant)

	//z = z // 4a
	c := floorDivision(z, new(big.Int).Mul(a, big.NewInt(4)))

	return NewClassGroup(a, b, c)
}

func NewClassGroupFromBytesDiscriminant(buf []byte, discriminant *big.Int) (*ClassGroup, bool) {

	int_size_bits := discriminant.BitLen()

	//add additional one byte for sign
	int_size := (int_size_bits + 16) >> 4

	//make sure the input byte buffer size matches with discriminant's
	if len(buf) != int_size*2 {
		return nil, false
	}

	a := new(big.Int)

	//if a is a negative number
	if buf[0]&0x8 != 0 {
		two_s_complement_encoding(buf[:int_size], int_size)
		a.SetBytes(buf[1:int_size])
		a = new(big.Int).Neg(a)
	} else {
		a.SetBytes(buf[1:int_size])
	}

	b := new(big.Int)

	//if a is a negative number
	if buf[int_size] != 0 {
		two_s_complement_encoding(buf[int_size:], int_size)
		b.SetBytes(buf[int_size+1 : 2*int_size])
		b = new(big.Int).Neg(b)
	} else {
		b.SetBytes(buf[int_size+1 : 2*int_size])
	}

	return NewClassGroupFromAbDiscriminant(a, b, discriminant), true
}

func IdentityForDiscriminant(d *big.Int) *ClassGroup {
	return NewClassGroupFromAbDiscriminant(big.NewInt(1), big.NewInt(1), d)
}

func (group *ClassGroup) Normalized() *ClassGroup {
	a := new(big.Int).Set(group.a)
	b := new(big.Int).Set(group.b)
	c := new(big.Int).Set(group.c)

	//if b > -a && b <= a:
	if (b.Cmp(new(big.Int).Neg(a)) == 1) && (b.Cmp(a) < 1) {
		return group
	}

	//r = (a - b) // (2 * a)
	r := new(big.Int).Sub(a, b)
	r = floorDivision(r, new(big.Int).Mul(a, big.NewInt(2)))

	//b, c = b + 2 * r * a, a * r * r + b * r + c
	t := new(big.Int).Mul(big.NewInt(2), r)
	t.Mul(t, a)
	oldB := new(big.Int).Set(b)
	b.Add(b, t)

	x := new(big.Int).Mul(a, r)
	x.Mul(x, r)
	y := new(big.Int).Mul(oldB, r)
	c.Add(c, x)
	c.Add(c, y)

	return NewClassGroup(a, b, c)
}

func (group *ClassGroup) Reduced() *ClassGroup {
	g := group.Normalized()
	a := new(big.Int).Set(g.a)
	b := new(big.Int).Set(g.b)
	c := new(big.Int).Set(g.c)

	//while a > c or (a == c and b < 0):
	for (a.Cmp(c) == 1) || ((a.Cmp(c) == 0) && (b.Sign() == -1)) {
		//s = (c + b) // (c + c)
		s := new(big.Int).Add(c, b)
		s = floorDivision(s, new(big.Int).Add(c, c))

		//a, b, c = c, -b + 2 * s * c, c * s * s - b * s + a
		oldA := new(big.Int).Set(a)
		oldB := new(big.Int).Set(b)
		a = new(big.Int).Set(c)

		b.Neg(b)
		x := new(big.Int).Mul(big.NewInt(2), s)
		x.Mul(x, c)
		b.Add(b, x)

		c.Mul(c, s)
		c.Mul(c, s)
		oldB.Mul(oldB, s)
		c.Sub(c, oldB)
		c.Add(c, oldA)
	}

	return NewClassGroup(a, b, c).Normalized()
}

func (group *ClassGroup) identity() *ClassGroup {
	return NewClassGroupFromAbDiscriminant(big.NewInt(1), big.NewInt(1), group.Discriminant())
}

func (group *ClassGroup) Discriminant() *big.Int {
	if group.d == nil {
		d := new(big.Int).Set(group.b)
		d.Mul(d, d)
		a := new(big.Int).Set(group.a)
		a.Mul(a, group.c)
		a.Mul(a, big.NewInt(4))
		d.Sub(d, a)

		group.d = d
	}
	return group.d
}

func (group *ClassGroup) Multiply(other *ClassGroup) *ClassGroup {
	//a1, b1, c1 = self.reduced()
	x := group.Reduced()

	//a2, b2, c2 = other.reduced()
	y := other.Reduced()

	//g = (b2 + b1) // 2
	g := new(big.Int).Add(x.b, y.b)
	g = floorDivision(g, big.NewInt(2))

	//h = (b2 - b1) // 2
	h := new(big.Int).Sub(y.b, x.b)
	h = floorDivision(h, big.NewInt(2))

	//w = mod.gcd(a1, a2, g)
	w1 := allInputValueGCD(y.a, g)
	w := allInputValueGCD(x.a, w1)

	//j = w
	j := new(big.Int).Set(w)
	//r = 0
	r := big.NewInt(0)
	//s = a1 // w
	s := floorDivision(x.a, w)
	//t = a2 // w
	t := floorDivision(y.a, w)
	//u = g // w
	u := floorDivision(g, w)

	/*
		solve these equations for k, l, m

		k * t - l * s = h
		k * u - m * s = c2
		l * u - m * t = c1

		solve
		(tu)k - (hu + sc) = 0 mod st
		k = (- hu - sc) * (tu)^-1
	*/

	//k_temp, constant_factor = mod.solve_mod(t * u, h * u + s * c1, s * t)
	b := new(big.Int).Mul(h, u)
	sc := new(big.Int).Mul(s, x.c)
	b.Add(b, sc)
	k_temp, constant_factor := SolveMod(new(big.Int).Mul(t, u), b, new(big.Int).Mul(s, t))

	//n, constant_factor_2 = mod.solve_mod(t * constant_factor, h - t * k_temp, s)
	n, _ := SolveMod(new(big.Int).Mul(t, constant_factor), new(big.Int).Sub(h, new(big.Int).Mul(t, k_temp)), s)

	//k = k_temp + constant_factor * n
	k := new(big.Int).Add(k_temp, new(big.Int).Mul(constant_factor, n))

	//l = (t * k - h) // s
	l := floorDivision(new(big.Int).Sub(new(big.Int).Mul(t, k), h), s)

	//m = (t * u * k - h * u - s * c1) // (s * t)
	tuk := new(big.Int).Mul(t, u)
	tuk.Mul(tuk, k)

	hu := new(big.Int).Mul(h, u)

	tuk.Sub(tuk, hu)
	tuk.Sub(tuk, sc)

	st := new(big.Int).Mul(s, t)
	m := floorDivision(tuk, st)

	//a3 = s * t - r * u
	ru := new(big.Int).Mul(r, u)
	a3 := st.Sub(st, ru)

	//b3 = (j * u + m * r) - (k * t + l * s)
	ju := new(big.Int).Mul(j, u)
	mr := new(big.Int).Mul(m, r)
	ju = ju.Add(ju, mr)

	kt := new(big.Int).Mul(k, t)
	ls := new(big.Int).Mul(l, s)
	kt = kt.Add(kt, ls)

	b3 := ju.Sub(ju, kt)

	//c3 = k * l - j * m
	kl := new(big.Int).Mul(k, l)
	jm := new(big.Int).Mul(j, m)

	c3 := kl.Sub(kl, jm)
	return NewClassGroup(a3, b3, c3).Reduced()
}

func (group *ClassGroup) Pow(n int64) *ClassGroup {
	x := CloneClassGroup(group)
	items_prod := group.identity()

	for n > 0 {
		if n&1 == 1 {
			items_prod = items_prod.Multiply(x)
		}
		x = x.Square()
		n >>= 1
	}
	return items_prod
}

func (group *ClassGroup) BigPow(n *big.Int) *ClassGroup {
	x := CloneClassGroup(group)
	items_prod := group.identity()

	p := new(big.Int).Set(n)
	for p.Sign() > 0 {
		if p.Bit(0) == 1 {
			items_prod = items_prod.Multiply(x)
		}
		x = x.Square()
		p.Rsh(p, 1)
	}
	return items_prod
}

func (group *ClassGroup) Square() *ClassGroup {

	//Solve bk ≡ c mod a, t
	//the solutions to which have the form k = µ + νn for all n ∈ Z.
	u, _ := SolveMod(group.b, group.c, group.a)
	//aa := fmt.Sprintf("%d",u)
	//fmt.Print(aa)

	//A = a
	A := new(big.Int).Mul(group.a, group.a)

	//B = b − 2aµ,
	au := new(big.Int).Mul(group.a, u)
	B := new(big.Int).Sub(group.b, new(big.Int).Mul(au, big.NewInt(2)))

	//C = µ ^ 2 - (bµ−c)//a
	C := new(big.Int).Mul(u, u)
	m := new(big.Int).Mul(group.b, u)
	m = new(big.Int).Sub(m, group.c)
	m = floorDivision(m, group.a)
	C = new(big.Int).Sub(C, m)

	return NewClassGroup(A, B, C).Reduced()
}

func (group *ClassGroup) SquareUsingMultiply() *ClassGroup {
	//a1, b1, c1 = self.reduced()
	x := group.Reduced()

	//g = b1
	g := x.b
	//h = 0
	h := big.NewInt(0)

	//w = mod.gcd(a1, g)
	w := allInputValueGCD(x.a, g)

	//j = w
	j := new(big.Int).Set(w)
	//r = 0
	r := big.NewInt(0)
	//s = a1 // w
	s := floorDivision(x.a, w)
	//t = s
	t := s
	//u = g // w
	u := floorDivision(g, w)

	/*
		solve these equations for k, l, m

		k * t - l * s = h
		k * u - m * s = c2
		l * u - m * t = c1

		solve
		(tu)k - (hu + sc) = 0 mod st
		k = (- hu - sc) * (tu)^-1
	*/

	//k_temp, constant_factor = mod.solve_mod(t * u, h * u + s * c1, s * t)
	b := new(big.Int).Mul(h, u)
	sc := new(big.Int).Mul(s, x.c)
	b.Add(b, sc)
	k_temp, constant_factor := SolveMod(new(big.Int).Mul(t, u), b, new(big.Int).Mul(s, t))

	//n, constant_factor_2 = mod.solve_mod(t * constant_factor, h - t * k_temp, s)
	n, _ := SolveMod(new(big.Int).Mul(t, constant_factor), new(big.Int).Sub(h, new(big.Int).Mul(t, k_temp)), s)

	//k = k_temp + constant_factor * n
	k := new(big.Int).Add(k_temp, new(big.Int).Mul(constant_factor, n))

	//l = (t * k - h) // s
	l := floorDivision(new(big.Int).Sub(new(big.Int).Mul(t, k), h), s)

	//m = (t * u * k - h * u - s * c1) // (s * t)
	tuk := new(big.Int).Mul(t, u)
	tuk.Mul(tuk, k)

	hu := new(big.Int).Mul(h, u)

	tuk.Sub(tuk, hu)
	tuk.Sub(tuk, sc)

	st := new(big.Int).Mul(s, t)
	m := floorDivision(tuk, st)

	//a3 = s * t - r * u
	ru := new(big.Int).Mul(r, u)
	a3 := st.Sub(st, ru)

	//b3 = (j * u + m * r) - (k * t + l * s)
	ju := new(big.Int).Mul(j, u)
	mr := new(big.Int).Mul(m, r)
	ju = ju.Add(ju, mr)

	kt := new(big.Int).Mul(k, t)
	ls := new(big.Int).Mul(l, s)
	kt = kt.Add(kt, ls)

	b3 := ju.Sub(ju, kt)

	//c3 = k * l - j * m
	kl := new(big.Int).Mul(k, l)
	jm := new(big.Int).Mul(j, m)

	c3 := kl.Sub(kl, jm)

	return NewClassGroup(a3, b3, c3).Reduced()
}

//encoding for a, b based on discriminant's size
//encoding using two's complement for a, b
//the first byte is sign, if a < 0, buf[0] = 0xff else buf[0] = 0
func (group *ClassGroup) Serialize() []byte {
	r := group.Reduced()
	int_size_bits := group.Discriminant().BitLen()
	int_size := (int_size_bits + 16) >> 4

	buf := make([]byte, int_size)
	a_bytes := r.a.Bytes()
	copy(buf[int_size-len(a_bytes):], a_bytes)

	//encode the negative number
	if r.a.Sign() == -1 {
		two_s_complement_encoding(buf, len(a_bytes))
	}

	buf2 := make([]byte, int_size)
	b_bytes := r.b.Bytes()
	copy(buf2[int_size-len(b_bytes):], b_bytes)

	//encode the negative number
	if r.b.Sign() == -1 {
		two_s_complement_encoding(buf2, len(b_bytes))
	}

	return append(buf, buf2...)
}

//encoding for a, b based on discriminant's size
//the first byte is sign, if a < 0, buf[0] = 1 else buf[0] = 0
func (group *ClassGroup) Serialize2() []byte {
	r := group.Reduced()
	int_size_bits := group.Discriminant().BitLen()

	//add additional one byte for sign
	int_size := (int_size_bits + 16) >> 4

	buf := make([]byte, int_size)
	a_bytes := r.a.Bytes()
	copy(buf[int_size-len(a_bytes):], a_bytes)

	//encode the negative number
	if r.a.Sign() == -1 {
		two_s_complement_encoding(buf, len(a_bytes))
	}

	buf2 := make([]byte, int_size)
	b_bytes := r.b.Bytes()
	copy(buf2[int_size-len(b_bytes):], b_bytes)

	//encode the negative number
	if r.b.Sign() == -1 {
		two_s_complement_encoding(buf2, len(b_bytes))
	}

	return append(buf, buf2...)
}

func (group *ClassGroup) Equal(other *ClassGroup) bool {
	g := group.Reduced()
	o := other.Reduced()

	return (g.a.Cmp(o.a) == 0 && g.b.Cmp(o.b) == 0 && g.c.Cmp(o.c) == 0)
}
