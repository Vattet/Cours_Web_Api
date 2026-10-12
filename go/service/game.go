package service

import (
	"context"

	"rest-api/model"
)

const defaultLimit = 10

type GameReader interface {
	FindAll(ctx context.Context, limit int64) ([]model.Game, error)
}

type GameService struct {
	reader	GameReader
}

func NewGameService(reader GameReader) *GameService {
	return &GameService{reader: reader}
}

func (s *GameService) List(ctx context.Context) ([]model.Game, error) {
	return s.reader.FindAll(ctx, defaultLimit)
}
