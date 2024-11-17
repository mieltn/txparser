package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/mieltn/txparser/internal/domain"
	"github.com/mieltn/txparser/internal/logger"
	"github.com/mieltn/txparser/pkg/eth"
	"io"
	"math/big"
	"net/http"
	"time"
)

const (
	contentType = "application/json"

	rpcVersion          = "2.0"
	methodBlockNumber   = "eth_blockNumber"
	methodBlockByNumber = "eth_getBlockByNumber"

	requestId = 0
)

type client struct {
	l       logger.Logger
	cl      *http.Client
	url     string
	retry   int
	retryIn int
}

func New(l logger.Logger, url string, retry, retryIn, timeoutSec int) *client {
	cl := http.DefaultClient
	cl.Timeout = time.Second * time.Duration(timeoutSec)
	return &client{
		l:       l,
		cl:      cl,
		url:     url,
		retry:   retry,
		retryIn: retryIn,
	}
}

func (c *client) GetBlockNumber(ctx context.Context) (*big.Int, error) {
	body := map[string]any{
		"jsonrpc": rpcVersion,
		"method":  methodBlockNumber,
		"params":  []map[string]any{},
	}
	raw, err := c.do(ctx, body)
	if err != nil {
		return nil, errors.Join(domain.ErrEthClient, err)
	}

	var resp EthBlockNumberResponse
	if err = json.Unmarshal(raw, &resp); err != nil {
		return nil, errors.Join(domain.ErrEthClient, err)
	}

	number, err := eth.ParseBigInt256(resp.Result)
	if err != nil {
		return nil, errors.Join(domain.ErrEthClient, err)
	}

	return number, nil
}

func (c *client) GetBlockByNumber(ctx context.Context, blockNumber string) ([]domain.Transaction, error) {
	body := map[string]any{
		"jsonrpc": rpcVersion,
		"method":  methodBlockByNumber,
		"params": []any{
			blockNumber,
			true,
		},
	}

	raw, err := c.do(ctx, body)
	if err != nil {
		return nil, errors.Join(domain.ErrEthClient, err)
	}

	var resp EthGetBlockByNumberResponse
	if err = json.Unmarshal(raw, &resp); err != nil {
		return nil, errors.Join(domain.ErrEthClient, err)
	}

	var txs []domain.Transaction
	for _, item := range resp.Result.Transactions {
		tx, err := item.toDomain()
		if err != nil {
			c.l.Errorf("failed to convert tx: hash %s, err: %v", item.Hash, errors.Join(domain.ErrEthClient, err))
			continue
		}
		txs = append(txs, tx)
	}

	return txs, nil

}

// TODO: implement retries
// TODO: add rate limiter
func (c *client) do(ctx context.Context, body map[string]any) ([]byte, error) {
	body["id"] = requestId
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := c.cl.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return raw, nil
}
