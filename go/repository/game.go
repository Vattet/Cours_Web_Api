package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"rest-api/model"
)

type GameRepository struct {
	collection	*mongo.Collection
}

func NewGameRepository(collection *mongo.Collection) *GameRepository {
	return &GameRepository{collection: collection}
}

func (r *GameRepository) FindAll(ctx context.Context, limit int64) ([]model.Game, error) {
	cursor, err := r.collection.Find(ctx, bson.D{}, options.Find().SetLimit(limit))
	if err != nil {
		return nil, err
	}
	var games []model.Game
	if err := cursor.All(ctx, &games); err != nil {
		return nil, err
	}
	return games, nil
}
