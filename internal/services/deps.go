package services

import (
	"context"
	"github.com/mieltn/txparser/internal/domain"
	"math/big"
)

type (
	addressesRepo interface {
		Create(ctx context.Context, address string) error
		IsSubscribed(ctx context.Context, address string) bool
	}
	transactionsRepo interface {
		Create(ctx context.Context, addr string, tx domain.Transaction) (int64, error)
		ByAddress(ctx context.Context, addr string) ([]domain.Transaction, error)
	}
	client interface {
		GetBlockNumber(ctx context.Context) (*big.Int, error)
		GetBlockByNumber(ctx context.Context, blockNumber string) ([]domain.Transaction, error)
	}
	lastProcessedBlockNumberRepo interface {
		Get(ctx context.Context) (*big.Int, error)
		Set(ctx context.Context, lastProcessedBlockNumber *big.Int) error
	}
)
