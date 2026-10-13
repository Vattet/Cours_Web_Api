package service

import "rest-api/model"

const defaultLimit = 10

type GameReader interface {
	FindAll(limit int64) ([]model.Game, error)
}

func ListGames(r GameReader) ([]model.Game, error) {
	return r.FindAll(defaultLimit)
}
