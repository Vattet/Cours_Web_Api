package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"rest-api/model"
	"rest-api/service"
)

func parseInt64(raw string) (int64, error) {
	if raw == "" {
		return 0, nil
	}
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || value < 0 {
		return 0, errors.New("invalid value")
	}
	return value, nil
}

func GetGames(reader service.GameReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		params := map[string]string{}
		for key, values := range query {
			if len(values) > 0 {
				params[key] = values[0]
			}
		}

		limit, err := parseInt64(query.Get("limit"))
		if err != nil {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		skip, err := parseInt64(query.Get("skip"))
		if err != nil {
			http.Error(w, "invalid skip", http.StatusBadRequest)
			return
		}

		page, err := service.ListGames(reader, params, limit, skip)

		if errors.Is(err, model.ErrInvalidFilter) {
			http.Error(w, "invalid filter value", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(page)
	}
}
