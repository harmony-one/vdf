package vdf_go

import "math/big"

//Floor Division for big.Int
//Reference : Division and Modulus for Computer Scientists
//https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/divmodnote.pdf
//Golang only has Euclid division and T-division
func floorDivision(x, y *big.Int) *big.Int {
	var r big.Int
	q, _ := new(big.Int).QuoRem(x, y, &r)

	if (r.Sign() == 1 && y.Sign() == -1) || (r.Sign() == -1 && y.Sign() == 1) {
		q.Sub(q, big.NewInt(1))
	}

	return q
}
