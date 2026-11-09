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

func GetGameByID(gameReader service.GameReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}

		game, err := service.GetGame(r.Context(), gameReader, id)
		if errors.Is(err, model.ErrGameNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if err != nil {
			log.Printf("GetGameByID %d: %v", id, err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(game)
	}
}

const maxBodySize = 1 << 20

func PostGame(writer service.GameWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

		var in model.CreateGame
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&in); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body: "+err.Error())
			return
		}

		game, err := service.CreateGame(r.Context(), writer, in)
		switch {
		case errors.Is(err, model.ErrInvalidGame):
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		case errors.Is(err, model.ErrGameConflict):
			writeError(w, http.StatusConflict, err.Error())
			return
		case err != nil:
			log.Printf("PostGame (user %s): %v", middleware.User(r.Context()), err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", fmt.Sprintf("/api/v1/games/%d", game.ID))
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(game)
	}
}

func PutGame(writer service.GameWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
		var in model.CreateGame
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&in); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body: "+err.Error())
			return
		}

		game, err := service.UpdateGame(r.Context(), writer, id, in)
		if errors.Is(err, model.ErrInvalidGame) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		if errors.Is(err, model.ErrGameNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if err != nil {
			log.Printf("PutGame %d (user %s): %v", id, middleware.User(r.Context()), err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(game)
	}
}

func DeleteGame(writer service.GameWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}

		err = service.DeleteGame(r.Context(), writer, id)
		if errors.Is(err, model.ErrGameNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if err != nil {
			log.Printf("DeleteGame %d (user %s): %v", id, middleware.User(r.Context()), err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
