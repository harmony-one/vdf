package main

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"math/big"
	"vdf_go"
)

func TestClassDiscriminant(t *testing.T) {
	t12_11_3 := vdf_go.NewClassGroup(big.NewInt(12), big.NewInt(11), big.NewInt(3))
	assert.Equal(t, "-23", t12_11_3.Discriminant().String(), "they should be equal")

	t93_109_32 := vdf_go.NewClassGroup(big.NewInt(93), big.NewInt(109), big.NewInt(32))
	assert.Equal(t, "-23", t93_109_32.Discriminant().String(), "they should be equal")

	D := big.NewInt(-103)
	e_id := vdf_go.IdentityForDiscriminant(D)
	assert.Equal(t, vdf_go.NewClassGroup(big.NewInt(1), big.NewInt(1), big.NewInt(26)), e_id, "they should be equal")
	assert.Equal(t, e_id.Discriminant(), D, "they should be equal")

	e := vdf_go.NewClassGroupFromAbDiscriminant(big.NewInt(2), big.NewInt(1), D)
	assert.Equal(t, vdf_go.NewClassGroup(big.NewInt(2), big.NewInt(1), big.NewInt(13)), e, "they should be equal")
	assert.Equal(t, e.Discriminant(), D, "they should be equal")
}

func TestNormalized(t *testing.T) {
	f := vdf_go.NewClassGroup(big.NewInt(195751), big.NewInt(1212121), big.NewInt(1876411))
	assert.Equal(t, vdf_go.NewClassGroup(big.NewInt(195751), big.NewInt(37615), big.NewInt(1807)), f.Normalized(), "they should be equal")
}

func TestReduced(t *testing.T) {
	f := vdf_go.NewClassGroup(big.NewInt(195751), big.NewInt(1212121), big.NewInt(1876411))
	assert.Equal(t, vdf_go.NewClassGroup(big.NewInt(1), big.NewInt(1), big.NewInt(1)), f.Reduced(), "they should be equal")
}

func check(a, b, c *big.Int, t *testing.T) {
	r, s := vdf_go.SolveMod(a, b, c)
	b.Mod(b, c)

	for k := 0; k < 50; k++ {
		//a_coefficient = r + s * k
		a_coefficient := new(big.Int).Add(r, new(big.Int).Mul(s, big.NewInt(int64(k))))
		aac := a_coefficient.Mul(a_coefficient, a)
		aac.Mod(aac, c)
		assert.Equal(t, aac, b, fmt.Sprintf("diff when k = %d", k))
	}
}

func TestSolveMod(t *testing.T) {

	check(big.NewInt(3), big.NewInt(4), big.NewInt(5), t)
	check(big.NewInt(6), big.NewInt(8), big.NewInt(10), t)
	check(big.NewInt(12), big.NewInt(30), big.NewInt(7), t)
	check(big.NewInt(6), big.NewInt(15), big.NewInt(411), t)
	check(big.NewInt(192), big.NewInt(193), big.NewInt(863), t)
	check(big.NewInt(-565721958), big.NewInt(740), big.NewInt(4486780496), t)
	check(big.NewInt(565721958), big.NewInt(740), big.NewInt(4486780496), t)
	check(big.NewInt(-565721958), big.NewInt(-740), big.NewInt(4486780496), t)
	check(big.NewInt(565721958), big.NewInt(-740), big.NewInt(4486780496), t)
}

func TestMultiplication1(t *testing.T) {
	t12_11_3 := vdf_go.NewClassGroup(big.NewInt(12), big.NewInt(11), big.NewInt(3))
	t93_109_32 := vdf_go.NewClassGroup(big.NewInt(93), big.NewInt(109), big.NewInt(32))

	a := t12_11_3.Multiply(t93_109_32)
	assert.Equal(t, a, vdf_go.NewClassGroup(big.NewInt(1), big.NewInt(1), big.NewInt(6)), "they should be equal")
}

func TestMultiplication2(t *testing.T) {
	t12_11_3 := vdf_go.NewClassGroup(big.NewInt(12), big.NewInt(11), big.NewInt(3))
	t93_109_32 := vdf_go.NewClassGroup(big.NewInt(93), big.NewInt(109), big.NewInt(32))

	x := vdf_go.CloneClassGroup(t12_11_3)
	y := t12_11_3.Multiply(x)
	assert.Equal(t, y, vdf_go.NewClassGroup(big.NewInt(2), big.NewInt(1), big.NewInt(3)), "they should be equal")

	x = vdf_go.CloneClassGroup(t93_109_32)
	y = t93_109_32.Multiply(x)
	assert.Equal(t, y, vdf_go.NewClassGroup(big.NewInt(2), big.NewInt(-1), big.NewInt(3)), "they should be equal")
}

func TestMultiplication3(t *testing.T) {
	t12_11_3 := vdf_go.NewClassGroup(big.NewInt(12), big.NewInt(11), big.NewInt(3))
	t93_109_32 := vdf_go.NewClassGroup(big.NewInt(93), big.NewInt(109), big.NewInt(32))

	a := t12_11_3.Multiply(t93_109_32)
	assert.Equal(t, a, vdf_go.NewClassGroup(big.NewInt(1), big.NewInt(1), big.NewInt(6)), "they should be equal")
}

func TestMultiplication4(t *testing.T) {
	x := vdf_go.NewClassGroup(big.NewInt(-565721958), big.NewInt(-740), big.NewInt(4486780496))
	y := vdf_go.NewClassGroup(big.NewInt(565721958), big.NewInt(740), big.NewInt(4486780496))

	a := x.Multiply(y)
	assert.Equal(t, a, vdf_go.NewClassGroup(big.NewInt(-1), big.NewInt(0), big.NewInt(2538270247313468068)), "they should be equal")
}

func TestSquare1(t *testing.T) {
	x := vdf_go.NewClassGroup(big.NewInt(12), big.NewInt(11), big.NewInt(3))
	y := x.Square()
	assert.Equal(t, y, vdf_go.NewClassGroup(big.NewInt(2), big.NewInt(1), big.NewInt(3)), "they should be equal")
}

func TestSquare2(t *testing.T) {
	x := vdf_go.NewClassGroup(big.NewInt(93), big.NewInt(109), big.NewInt(32))
	y := x.Square()
	assert.Equal(t, y, vdf_go.NewClassGroup(big.NewInt(2), big.NewInt(-1), big.NewInt(3)), "they should be equal")
}

func TestSerialize(t *testing.T) {
	x := vdf_go.NewClassGroup(big.NewInt(-565721958), big.NewInt(-740), big.NewInt(4486780496))
	s := fmt.Sprintf("%02x", x.Serialize())
	assert.Equal(t, s, "ffde47c49afffffffd1c", "they should be equal")
}

func TestSerialize1(t *testing.T) {
	x := vdf_go.NewClassGroup(big.NewInt(-0x10000), big.NewInt(-740), big.NewInt(4486780496))
	s := fmt.Sprintf("%02x", x.Serialize())
	assert.Equal(t, s, "ffff0000fffffd1c", "they should be equal")
}

func TestDeSerialize1(t *testing.T) {
	str := "ff100000ffffffff"
	buf, _ := hex.DecodeString(str)
	x, _ := vdf_go.NewClassGroupFromBytesDiscriminant(buf, big.NewInt(4486780496111111))
	assert.Equal(t, str, hex.EncodeToString(x.Serialize()), "they should be equal")
}
