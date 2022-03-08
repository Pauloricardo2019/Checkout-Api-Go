package user

import (
	"context"
	"ravxcheckout/src/adapter/database/mongodb"
	"ravxcheckout/src/internal/model"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection
var ctx context.Context

func init() {
	collection, ctx = mongodb.GetCollection("user")
}

func GetAll() ([]model.User, error) {
	findResults, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var results []model.User
	err = findResults.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func GetByID(ID string) (*model.User, error) {
	filter := bson.D{{"_id", ID}}
	result := collection.FindOne(ctx, filter)

	user := &model.User{}
	err := result.Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func Create(user *model.User) error {
	user.ID = uuid.NewString()

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func Update(user *model.User) error {
	_, err := collection.UpdateOne(
		ctx,
		bson.D{{"_id", user.ID}},
		bson.D{{"$set", user}},
	)
	if err != nil {
		return err
	}

	return nil
}

func Delete(ID string) error {
	_, err := collection.DeleteOne(
		ctx,
		bson.D{{"_id", ID}},
	)
	if err != nil {
		return err
	}

	return nil
}
