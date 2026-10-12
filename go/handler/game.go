package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"rest-api/model"
)

type GameLister interface {
	List(ctx context.Context) ([]model.Game, error)
}

type GameHandler struct {
	service	GameLister
}

func NewGameHandler(service GameLister) *GameHandler {
	return &GameHandler{service: service}
}

func (h *GameHandler) List(w http.ResponseWriter, r *http.Request) {
	games, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}
