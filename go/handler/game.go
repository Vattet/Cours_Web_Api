package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"rest-api/middleware"
	"rest-api/model"
	"rest-api/service"
	"strconv"
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
			fmt.Println(values)
			for _, value := range values {
				params[key] = value
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

		if errors.Is(err, model.ErrInvalidFilter) { // == if(err == model.InvalidFilter)
			http.Error(w, "invalid filter key", http.StatusBadRequest)
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

const maxBodySize = 1 << 20 // 1 Mo
func PostGame(writer service.GameWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize) // protège contre les corps énormes

		var in model.CreateGame
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields() // refuse {"_id": 5} ou les fautes de frappe
		if err := decoder.Decode(&in); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body: "+err.Error())
			return
		}

		game, err := service.CreateGame(r.Context(), writer, in)
		switch {
		case errors.Is(err, model.ErrInvalidGame):
			writeError(w, http.StatusUnprocessableEntity, err.Error()) // 422
			return
		case errors.Is(err, model.ErrGameConflict):
			writeError(w, http.StatusConflict, err.Error()) // 409
			return
		case err != nil:
			log.Printf("PostGame (user %s): %v", middleware.User(r.Context()), err)
			writeError(w, http.StatusInternalServerError, "internal error") // ne jamais exposer l'erreur Mongo
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", fmt.Sprintf("/api/v1/games/%d", game.ID))
		w.WriteHeader(http.StatusCreated) // 201
		json.NewEncoder(w).Encode(game)
	}
}
