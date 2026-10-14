package service

import "rest-api/model"

const defaultLimit = 10

type GameReader interface {
	FindAll(limit int64, params map[string]string) ([]model.Game, error)
}

func ListGames(r GameReader, params map[string]string) ([]model.Game, error) {
	return r.FindAll(defaultLimit, params)
}
