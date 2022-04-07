package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang-stepik-2022q1/reditclone/config"
	"golang-stepik-2022q1/reditclone/pkg/log"
	"time"
)

func NewMongo() *mongo.Client {
	uri := fmt.Sprintf("mongodb://%s:%s", config.Cfg.MongoHost, config.Cfg.MongoPort)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Error("cant connect to mongo", log.Fields{"error": err.Error()})
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Error("cant ping mongo", log.Fields{"error": err.Error()})
	}
	return client
}
