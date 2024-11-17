package eth

import (
	"errors"
	"math/big"
)

var (
	errParseBigInt256 = errors.New("failed to parse big int256")
)

func ParseBigInt256(s string) (*big.Int, error) {
	value := new(big.Int)
	if s[2:] == "" {
		return value, nil
	}
	value, ok := value.SetString(s[2:], 16)
	if !ok {
		return nil, errParseBigInt256
	}
	return value, nil
}
