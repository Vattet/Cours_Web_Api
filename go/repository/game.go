package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"rest-api/model"
)

type GameRepository struct {
	*mongo.Collection
}

func (r *GameRepository) FindAll(limit int64) ([]model.Game, error) {
	ctx := context.Background()
	cursor, err := r.Find(ctx, bson.D{}, options.Find().SetLimit(limit))
	if err != nil {
		return nil, err
	}
	var games []model.Game
	if err := cursor.All(ctx, &games); err != nil {
		return nil, err
	}
	return games, nil
}
