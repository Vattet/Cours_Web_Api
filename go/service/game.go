package service

import "rest-api/model"

const defaultLimit = 10

type GameReader interface {
	FindAll(limit int64, title string, year int) ([]model.Game, error)
}

func ListGames(r GameReader, title string, year int) ([]model.Game, error) {
	return r.FindAll(defaultLimit, title, year)
}
