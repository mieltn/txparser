package handlers

import (
	"context"
	"github.com/mieltn/txparser/internal/domain"
)

type txparserService interface {
	GetCurrentBlock(ctx context.Context) (uint64, error)
	Subscribe(ctx context.Context, address string) bool
	GetTransactions(ctx context.Context, address string) ([]domain.Transaction, error)
}
