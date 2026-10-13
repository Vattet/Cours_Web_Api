package router

import (
	"net/http"

	"rest-api/handler"
	"rest-api/service"
)

func Routes(reader service.GameReader) {
	http.HandleFunc("GET /api/v1/games", handler.GetGames(reader))
}
