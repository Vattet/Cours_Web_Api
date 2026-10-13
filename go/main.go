package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"rest-api/repository"
	"rest-api/router"
)

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func main() {
	if err := loadEnv(getenv("ENV_FILE", "../database/.env")); err != nil {
		log.Fatal(err)
	}

	uri := "mongodb://" + os.Getenv("MONGO_ROOT_USERNAME") + ":" + os.Getenv("MONGO_ROOT_PASSWORD") +
		"@" + getenv("MONGO_HOST", "localhost") + ":" + getenv("MONGO_PORT", "27017") + "/?authSource=admin"

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	collection := client.Database(os.Getenv("MONGO_DATABASE"))
	games := collection.Collection("games")
	router.Routes(&repository.GameRepository{games})

	log.Fatal(http.ListenAndServe(":"+getenv("PORT", "8080"), nil))
}
