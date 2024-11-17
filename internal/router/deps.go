package router

import "net/http"

type (
	addressesHandler interface {
		Subscribe(http.ResponseWriter, *http.Request)
	}
	transactionsHandler interface {
		GetCurrentBlock(http.ResponseWriter, *http.Request)
		GetTransactions(http.ResponseWriter, *http.Request)
	}
)
