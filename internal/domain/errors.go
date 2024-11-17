package domain

import "errors"

var (
	ErrLastProcessedBlockHolder = errors.New("last processed block holder error")
	ErrTxparser                 = errors.New("txparser error")
	ErrEthClient                = errors.New("eth client request error")
	ErrAddressNotFound          = errors.New("address not found")
)
