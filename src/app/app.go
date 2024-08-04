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

func NewAppContext() *AppContext {
	config := LoadConfig()
	httpClient := NewHTTPClient()
	mongoClient := NewMongoClient(config)

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

func NewMongoClient(config Config) *mongo.Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.MongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
