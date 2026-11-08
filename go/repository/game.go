package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"rest-api/model"
)

type GameRepository struct {
	*mongo.Collection
}

func (r *GameRepository) FindAll(limit, skip int64, params map[string]string) ([]model.Game, error) {
	ctx := context.Background()

	filter, sort, err := buildQuery(params)
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
	filter, _, err := buildQuery(params)
	if err != nil {
		return 0, err
	}
	return r.CountDocuments(context.Background(), filter)
}
func (r *GameRepository) Create(ctx context.Context, game model.Game) (model.Game, error) {
	id, err := r.nextID(ctx)
	if err != nil {
		return model.Game{}, err
	}
	game.ID = id

	if _, err := r.InsertOne(ctx, game); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return model.Game{}, model.ErrGameConflict
		}
		return model.Game{}, err
	}
	return game, nil
}

func (r *GameRepository) nextID(ctx context.Context) (int, error) {
	var last struct {
		ID int `bson:"_id"`
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}}).SetProjection(bson.D{{Key: "_id", Value: 1}})

	err := r.FindOne(ctx, bson.D{}, opts).Decode(&last)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return 1, nil
	}
	if err != nil {
		return 0, err
	}
	return last.ID + 1, nil
}

func (r *GameRepository) FindByID(ctx context.Context, id int) (model.Game, error) {
	var game model.Game

	err := r.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&game)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.Game{}, model.ErrGameNotFound
	}
	if err != nil {
		return model.Game{}, err
	}
	return game, nil
}
