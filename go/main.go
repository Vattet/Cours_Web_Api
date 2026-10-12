package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"rest-api/handler"
	"rest-api/repository"
	"rest-api/router"
	"rest-api/service"
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

	uri := (&url.URL{
		Scheme:		"mongodb",
		User:		url.UserPassword(os.Getenv("MONGO_ROOT_USERNAME"), os.Getenv("MONGO_ROOT_PASSWORD")),
		Host:		getenv("MONGO_HOST", "localhost") + ":" + getenv("MONGO_PORT", "27017"),
		Path:		"/",
		RawQuery:	"authSource=admin",
	}).String()

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

	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_COLLECTION"))
	games := handler.NewGameHandler(service.NewGameService(repository.NewGameRepository(collection)))
	router.Register(games)

	log.Fatal(http.ListenAndServe(":"+getenv("PORT", "8080"), nil))
}
