package mongodb

import (
	"context"
	"ravxcheckout/src/adapter/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var db *mongo.Database
var ctx context.Context

func init() {
	localCtx, _ := context.WithTimeout(context.Background(), 60*time.Second)

	cfg := config.GetConfig()
	mongoClient, err := mongo.Connect(localCtx, options.Client().ApplyURI(cfg.MongoDBConnString))
	if err != nil {
		panic(err.Error())
	}

	// Ping the primary
	if err := mongoClient.Ping(localCtx, readpref.Primary()); err != nil {
		panic(err)
	}

	db = mongoClient.Database("ravxcheckout")
	ctx = localCtx
}

func GetCollection(collectionName string) (*mongo.Collection, context.Context) {
	conn, localCtx := GetConnection()
	return conn.Collection(collectionName), localCtx
}

func GetConnection() (*mongo.Database, context.Context) {
	return db, ctx
}
