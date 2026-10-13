package handler

import (
	"encoding/json"
	"net/http"

	"rest-api/service"
)

func GetGames(reader service.GameReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		games, err := service.ListGames(reader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(games)
	}
}
