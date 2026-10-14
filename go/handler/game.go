package handler

import (
	"encoding/json"
	"net/http"

	"rest-api/service"
)

func GetGames(gameReader service.GameReader) http.HandlerFunc {
	return func(writer http.ResponseWriter, reader *http.Request) {
		games, err := service.ListGames(gameReader)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(games)
	}
}

func GetGamesId(gameReader service.GameReader) http.HandlerFunc {
	return func(writer http.ResponseWriter, reader *http.Request) {
		id, err := reader.PathValue("id")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		games, err := service.ListGames(gameReader)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(games)
	}
}


