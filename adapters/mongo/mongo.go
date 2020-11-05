// Package db includes mongoDB implementation of application interfaces
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DBNAME = "icfs"
const Timeout = 5

func NewMongo(conURI string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conURI))
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to mongoDB instance")
	}
	return client.Database(DBNAME), nil
}

type MongoCol struct {
	Col *mongo.Collection
}

func (mcl *MongoCol) InsertOne(object interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()
	res, err := mcl.Col.InsertOne(ctx, object)
	if err != nil {
		return "", errors.Wrap(err, "failed to insert user into mongoDB")
	}
	return fmt.Sprintf("%v", res.InsertedID), nil
}

func (mcl *MongoCol) FindOne(filter interface{}) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()
	return mcl.Col.FindOne(ctx, filter)
}

func (mcl *MongoCol) DeleteOne(filter interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()
	_, err := mcl.Col.DeleteOne(ctx, filter)
	return errors.Wrap(err, "failed to delete object")
}

func (mcl *MongoCol) UpdateOne(filter, update interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()
	_, err := mcl.Col.UpdateOne(ctx, filter, update)
	return errors.Wrap(err, "failed to update object")
}
