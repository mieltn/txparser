package handlers

import (
	"encoding/json"
	"github.com/mieltn/txparser/api"
	"net/http"
)

func writeJson(w http.ResponseWriter, obj any) {
	raw, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if _, err = w.Write(raw); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeJsonWithError(w http.ResponseWriter, err error) {
	resp := api.Error{
		Message: err.Error(),
	}
	raw, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if _, err = w.Write(raw); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
