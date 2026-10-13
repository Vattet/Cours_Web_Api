package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type HealthRepository struct {
	*mongo.Client
}

func (r *HealthRepository) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx, readpref.Primary())
}
