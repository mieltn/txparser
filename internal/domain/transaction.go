package domain

import "math/big"

type Transaction struct {
	Hash     string
	Amount   *big.Int
	FromAddr string
	ToAddr   string
}
