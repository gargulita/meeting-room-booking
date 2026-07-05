package handlers

import (
	"net/http"
)

type InfoHandler struct{}

func NewInfoHandler() *InfoHandler {
	return &InfoHandler{}
}

func (h *InfoHandler) Info(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
