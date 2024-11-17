package inmemory

import (
	"context"
	"github.com/mieltn/txparser/internal/logger"
	"math/big"
	"sync"
)

type processedBlockRepo struct {
	l               logger.Logger
	lastBlockNumber *big.Int
	mtx             *sync.Mutex
}

func NewProcessedBlockRepo(l logger.Logger, lastProcessedBlock *big.Int) *processedBlockRepo {
	return &processedBlockRepo{
		l:               l,
		lastBlockNumber: lastProcessedBlock,
		mtx:             &sync.Mutex{},
	}
}

func (r *processedBlockRepo) Get(ctx context.Context) (*big.Int, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return r.lastBlockNumber, nil
}

func (r *processedBlockRepo) Set(ctx context.Context, blockNumber *big.Int) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if r.lastBlockNumber.Cmp(blockNumber) < 0 {
		r.lastBlockNumber = blockNumber
	}
	return nil
}
