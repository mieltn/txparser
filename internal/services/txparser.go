package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/mieltn/txparser/internal/domain"
	"github.com/mieltn/txparser/internal/logger"
	"math/big"
	"strings"
	"sync"
)

type job struct {
	BlockNumber *big.Int
}

type txparser struct {
	l                            logger.Logger
	wg                           *sync.WaitGroup
	cl                           client
	transactionsRepo             transactionsRepo
	addressesRepo                addressesRepo
	lastProcessedBlockNumberRepo lastProcessedBlockNumberRepo
	jobs                         chan job
	nWorkers                     int
}

func NewTxparser(
	l logger.Logger,
	rpcClient client,
	transactionsRepo transactionsRepo,
	addressesRepo addressesRepo,
	lastProcessedBlockNumberRepo lastProcessedBlockNumberRepo,
	nWorkers int,
) *txparser {
	return &txparser{
		l:                            l,
		cl:                           rpcClient,
		transactionsRepo:             transactionsRepo,
		addressesRepo:                addressesRepo,
		lastProcessedBlockNumberRepo: lastProcessedBlockNumberRepo,
		jobs:                         make(chan job, 100),
		nWorkers:                     nWorkers,
	}
}

func (tp *txparser) Start(ctx context.Context) {
	tp.wg = &sync.WaitGroup{}
	for i := 0; i < tp.nWorkers; i++ {
		tp.wg.Add(1)
		go tp.queryChain(ctx)
	}
	tp.l.Infof("started tx parser with %d workers", tp.nWorkers)
}

func (tp *txparser) Stop() {
	tp.wg.Wait()
	close(tp.jobs)
	tp.l.Infof("stopped tx parser")
}

func (tp *txparser) ParseBlocks(ctx context.Context) error {
	lastProcessed, err := tp.lastProcessedBlockNumberRepo.Get(ctx)
	if err != nil {
		return errors.Join(err, domain.ErrLastProcessedBlockHolder)
	}

	lastBlock, err := tp.cl.GetBlockNumber(ctx)
	if err != nil {
		return errors.Join(err, domain.ErrTxparser)
	}

	for i := new(big.Int).Set(lastProcessed); i.Cmp(lastBlock) < 0; i.Add(i, big.NewInt(1)) {
		tp.jobs <- job{
			BlockNumber: new(big.Int).Set(i),
		}
	}

	return nil
}

func (tp *txparser) GetCurrentBlock(ctx context.Context) (uint64, error) {
	block, err := tp.lastProcessedBlockNumberRepo.Get(ctx)
	if err != nil {
		return 0, errors.Join(domain.ErrTxparser, err)
	}
	return block.Uint64(), nil
}

func (tp *txparser) Subscribe(ctx context.Context, address string) bool {
	if err := tp.addressesRepo.Create(ctx, strings.ToLower(address)); err != nil {
		tp.l.Errorf("failed to subscribe: %v", errors.Join(domain.ErrTxparser, err))
		return false
	}
	return true
}

func (tp *txparser) GetTransactions(ctx context.Context, address string) ([]domain.Transaction, error) {
	txs, err := tp.transactionsRepo.ByAddress(ctx, strings.ToLower(address))
	if err != nil {
		return nil, errors.Join(domain.ErrTxparser, err)
	}
	return txs, nil
}

func (tp *txparser) queryChain(ctx context.Context) {
	defer tp.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case j := <-tp.jobs:
			blockHex := fmt.Sprintf("0x%s", j.BlockNumber.Text(16))
			blockTxs, err := tp.cl.GetBlockByNumber(ctx, blockHex)
			if err != nil {
				tp.l.Errorf(
					"failed to get block: number %s, err: %v",
					blockHex,
					errors.Join(domain.ErrTxparser, err),
				)
				continue
			}

			for _, tx := range blockTxs {
				if tp.addressesRepo.IsSubscribed(ctx, tx.FromAddr) {
					_, err = tp.transactionsRepo.Create(ctx, tx.FromAddr, tx)
					if err != nil {
						tp.l.Errorf(
							"failed to add transaction: addr %s, hash %s, err: %v",
							tx.FromAddr, tx.Hash,
							errors.Join(domain.ErrTxparser, err),
						)
					} else {
						tp.l.Infof("added transaction for %s: %+v", tx.FromAddr, tx)
					}
				}
				if tp.addressesRepo.IsSubscribed(ctx, tx.ToAddr) {
					_, err = tp.transactionsRepo.Create(ctx, tx.ToAddr, tx)
					if err != nil {
						tp.l.Errorf(
							"failed to add transaction: addr %s, hash %s, err: %v",
							tx.ToAddr, tx.Hash,
							errors.Join(domain.ErrTxparser, err),
						)
					} else {
						tp.l.Infof("added transaction for %s: %+v", tx.ToAddr, tx)
					}
				}
			}

			if err = tp.lastProcessedBlockNumberRepo.Set(ctx, j.BlockNumber); err != nil {
				tp.l.Errorf("failed update last block holder: %v", err)
			}
		}
	}
}
