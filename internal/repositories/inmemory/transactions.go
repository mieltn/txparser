package inmemory

import (
	"context"
	"github.com/mieltn/txparser/internal/domain"
	"github.com/mieltn/txparser/internal/logger"
	"math/big"
	"sync"
)

type transaction struct {
	Hash     string
	Amount   *big.Int
	FromAddr string
	ToAddr   string
}

type transactionsRepository struct {
	l    logger.Logger
	data map[string][]transaction
	mtx  sync.RWMutex
}

func NewTransactions(l logger.Logger) *transactionsRepository {
	data := make(map[string][]transaction, 100)
	return &transactionsRepository{
		l:    l,
		data: data,
		mtx:  sync.RWMutex{},
	}
}

func (a *transactionsRepository) Create(
	ctx context.Context,
	addr string,
	tx domain.Transaction,
) (int64, error) {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	a.data[addr] = append(a.data[addr], transaction{
		Hash:     tx.Hash,
		Amount:   tx.Amount,
		FromAddr: tx.FromAddr,
		ToAddr:   tx.ToAddr,
	})
	return 0, nil
}

// TODO: add pagination
func (a *transactionsRepository) ByAddress(
	ctx context.Context, addr string,
) ([]domain.Transaction, error) {
	a.mtx.RLock()
	defer a.mtx.RUnlock()
	var res []domain.Transaction
	for _, tx := range a.data[addr] {
		res = append(res, tx.toDomain())
	}
	return res, nil
}

func (t *transaction) toDomain() domain.Transaction {
	return domain.Transaction{
		Hash:     t.Hash,
		Amount:   t.Amount,
		FromAddr: t.FromAddr,
		ToAddr:   t.ToAddr,
	}
}
