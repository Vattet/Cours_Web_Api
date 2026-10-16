package router

import (
	"net/http"

	"rest-api/handler"
	"rest-api/middleware"

	"rest-api/service"
)

func Routes(reader service.GameReader, writer service.GameWriter, pinger handler.Pinger) {
	http.HandleFunc("GET /api/v1/games", handler.GetGames(reader))
	http.Handle("POST /api/v1/games", middleware.RequireAuth(handler.PostGame(writer)))

	http.HandleFunc("GET /health", handler.GetHealth(pinger))

	http.HandleFunc("POST /api/v1/login", handler.UserConnection())
}
