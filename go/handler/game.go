package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"rest-api/service"
)

func GetGames(reader service.GameReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := strings.TrimSpace(r.URL.Query().Get("title"))

		year := 0
		if raw := r.URL.Query().Get("release_year"); raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err != nil {
				http.Error(w, "invalid release_year", http.StatusBadRequest)
				return
			}
			year = parsed
		}

		games, err := service.ListGames(reader, title, year)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(games) == 0 {
			http.Error(w, "no games found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(games)
	}
}
