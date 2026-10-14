package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"rest-api/model"
	"rest-api/service"
)

func GetGames(reader service.GameReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := map[string]string{}
		for key, values := range r.URL.Query() {
			if len(values) > 0 {
				params[key] = values[0]
			}
		}

		games, err := service.ListGames(reader, params)

		if errors.Is(err, model.ErrInvalidFilter) {
			http.Error(w, "invalid filter value", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(games)
	}
}
