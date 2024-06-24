package main

import (
	"math/big"
	"strconv"
)

func StringToBigInt(v string) (*big.Int, error) {
	v64, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, err
	}
	return big.NewInt(v64), nil
}
