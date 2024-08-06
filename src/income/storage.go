package income

import (
	"context"

	"github.com/rafaljusiak/daily-dashboard/app"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IncomeDocument struct {
	Date         string  `bson:"date"`
	TimeHours    float64 `bson:"time_hours,omitempty"`
	HourlyRate   float64 `bson:"hourly_rate,omitempty"`
	ExchangeRate float64 `bson:"exchange_rate,omitempty"`
}

func getCollection(appCtx *app.AppContext) *mongo.Collection {
	return appCtx.MongoClient.Database("daily_dashboard").Collection("income")
}

func InsertDocument(
	appCtx *app.AppContext,
	ctx context.Context,
	document IncomeDocument,
) error {
	collection := getCollection(appCtx)

	_, err := collection.InsertOne(ctx, document)
	if err != nil {
		return err
	}

	return nil
}

func UpdateDocument(
	appCtx *app.AppContext,
	ctx context.Context,
	document IncomeDocument,
) error {
	collection := getCollection(appCtx)

	filter := bson.M{"date": document.Date}
	update := bson.M{"$set": document}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func GetDocumentByDate(
	appCtx *app.AppContext,
	ctx context.Context,
	date string,
) (IncomeDocument, error) {
	collection := getCollection(appCtx)
	filter := bson.M{"date": date}

	var document IncomeDocument
	err := collection.FindOne(ctx, filter).Decode(&document)
	if err != nil {
		return IncomeDocument{}, err
	}

	return document, nil
}
