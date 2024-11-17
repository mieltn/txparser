package client

import (
	"github.com/mieltn/txparser/internal/domain"
	"github.com/mieltn/txparser/pkg/eth"
)

type EthBlockNumberResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
}

type EthBlockTx struct {
	Hash  string `json:"hash"`
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}

type EthGetBlockByNumberResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Transactions []EthBlockTx `json:"transactions"`
	} `json:"result"`
	Id    int64     `json:"id"`
	Error *EthError `json:"error,omitempty"`
}

type EthError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (e *EthBlockTx) toDomain() (domain.Transaction, error) {
	amount, err := eth.ParseBigInt256(e.Value)
	return domain.Transaction{
		Hash:     e.Hash,
		Amount:   amount,
		FromAddr: e.From,
		ToAddr:   e.To,
	}, err
}
