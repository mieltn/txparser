package handlers

import (
	"github.com/mieltn/txparser/api"
	"github.com/mieltn/txparser/internal/logger"
	"net/http"
)

type addressesHandler struct {
	l       logger.Logger
	service txparserService
}

func NewAddresses(l logger.Logger, service txparserService) *addressesHandler {
	return &addressesHandler{l: l, service: service}
}

func (h *addressesHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	addr := r.PathValue("address")
	ok := h.service.Subscribe(r.Context(), addr)
	w.WriteHeader(http.StatusOK)
	writeJson(w, api.SubscribeResponse{
		Ok: ok,
	})
}
