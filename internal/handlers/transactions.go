package handlers

import (
	"errors"
	"github.com/mieltn/txparser/api"
	"github.com/mieltn/txparser/internal/domain"
	"github.com/mieltn/txparser/internal/logger"
	"net/http"
)

type transactionsHandler struct {
	l       logger.Logger
	service txparserService
}

func NewTransactions(l logger.Logger, service txparserService) *transactionsHandler {
	return &transactionsHandler{
		l:       l,
		service: service,
	}
}

func (h *transactionsHandler) GetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	block, err := h.service.GetCurrentBlock(r.Context())
	if err != nil {
		h.l.Errorf("failed to get current block: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		writeJsonWithError(w, err)
	}
	w.WriteHeader(http.StatusOK)
	writeJson(w, api.CurrentBlockResponse{
		BlockNumber: block,
	})
}

func (h *transactionsHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	addr := r.PathValue("address")
	txs, err := h.service.GetTransactions(r.Context(), addr)
	if errors.Is(err, domain.ErrAddressNotFound) {
		w.WriteHeader(http.StatusNotFound)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	writeJson(w, toGetTransactionsResponse(txs))

}

func toGetTransactionsResponse(txs []domain.Transaction) api.GetTransactionsResponse {
	var respTxs []api.Transaction
	for _, tx := range txs {
		respTxs = append(respTxs, api.Transaction{
			Hash:   tx.Hash,
			Amount: tx.Amount.Uint64(),
			From:   tx.FromAddr,
			To:     tx.ToAddr,
		})
	}
	return api.GetTransactionsResponse{
		Transactions: respTxs,
	}
}
