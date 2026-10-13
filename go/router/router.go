package router

import (
	"net/http"

	"rest-api/handler"
	"rest-api/service"
)

func Routes(reader service.GameReader, pinger handler.Pinger) {
	http.HandleFunc("GET /api/v1/games", handler.GetGames(reader))

	http.HandleFunc("GET /health", handler.GetHealth(pinger))

}
