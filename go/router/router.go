package router

import (
	"net/http"

	"rest-api/handler"
)

func Register(games *handler.GameHandler) {
	http.HandleFunc("GET /api/v1/games", games.List)
}
