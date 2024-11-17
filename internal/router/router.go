package router

import (
	"net/http"
)

type router struct {
	mux *http.ServeMux
}

func New(addresses addressesHandler, transactions transactionsHandler) *router {
	mx := http.NewServeMux()
	mx.HandleFunc("GET /current_block", transactions.GetCurrentBlock)
	mx.HandleFunc("POST /subscribe/{address}", addresses.Subscribe)
	mx.HandleFunc("GET /transactions/{address}", transactions.GetTransactions)

	return &router{
		mux: mx,
	}
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
