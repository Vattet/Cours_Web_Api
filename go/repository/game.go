package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"rest-api/model"
)

type GameRepository struct {
	*mongo.Collection
}

func (r *GameRepository) FindAll(limit, skip int64, params map[string]string) ([]model.Game, error) {
	ctx := context.Background()

	filter, err := buildFilter(params)
	if err != nil {
		return nil, err
	}

	sort, err := buildSort(params)
	if err != nil {
		return nil, err
	}

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(sort)

	cursor, err := r.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var games []model.Game
	if err := cursor.All(ctx, &games); err != nil {
		return nil, err
	}

	if games == nil {
		games = []model.Game{}
	}

	return games, nil
}

func (r *GameRepository) Count(params map[string]string) (int64, error) {
	filter, err := buildFilter(params)
	if err != nil {
		return 0, nil
	}
	return r.CountDocuments(context.Background(), filter)
}
