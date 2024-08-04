package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AppContext struct {
	Config      Config
	HTTPClient  *http.Client
	MongoClient *mongo.Client
}

func NewAppContext(ctx context.Context) *AppContext {
	config := LoadConfig()
	httpClient := NewHTTPClient()
	mongoClient := NewMongoClient(ctx, config)

	return &AppContext{
		Config:      config,
		HTTPClient:  httpClient,
		MongoClient: mongoClient,
	}
}

func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}

func NewMongoClient(ctx context.Context, config Config) *mongo.Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.MongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
