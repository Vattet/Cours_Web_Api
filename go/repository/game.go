package repository

import (
	"context"
	"regexp"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"rest-api/model"
)

type GameRepository struct {
	*mongo.Collection
}

func (r *GameRepository) FindAll(limit int64, title string, year int) ([]model.Game, error) {
	ctx := context.Background()

	filter := bson.D{}
	if title != "" {
		filter = append(filter, bson.E{
			Key: "title",
			Value: bson.D{
				{Key: "$regex", Value: regexp.QuoteMeta(title)},
				{Key: "$options", Value: "i"}, // un <<changement>> entre gros <<>> pour que si I == i = true (change la regle de comparaison plutot)

			},
		})
	}
	if year != 0 {
		filter = append(filter, bson.E{Key: "release_year", Value: year})
	}
	cursor, err := r.Find(ctx, filter, options.Find().SetLimit(limit))
	if err != nil {
		return nil, err
	}
	var games []model.Game
	if err := cursor.All(ctx, &games); err != nil {
		return nil, err
	}
	return games, nil
}
