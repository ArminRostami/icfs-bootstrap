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

type mongoDB struct {
	db *mongo.Database
}

func NewMongo(conURI string) (*mongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conURI))
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to mongoDB instance")
	}
	return &mongoDB{db: client.Database(DBNAME)}, nil
}

func (mdb *mongoDB) InsertOne(collection string, object interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()
	res, err := mdb.db.Collection(collection).InsertOne(ctx, object)
	return fmt.Sprintf("%v", res.InsertedID), errors.Wrap(err, "failed to insert user into mongoDB")
}

func (mdb *mongoDB) FindOne(collection string, filter interface{}) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()
	return mdb.db.Collection(collection).FindOne(ctx, filter)
}

func (mdb *mongoDB) DeleteOne(collection string, filter interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()
	_, err := mdb.db.Collection(collection).DeleteOne(ctx, filter)
	return errors.Wrap(err, "failed to delete object")
}

func (mdb *mongoDB) UpdateOne(collection string, filter, update interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()
	_, err := mdb.db.Collection(collection).ReplaceOne(ctx, filter, update)
	return errors.Wrap(err, "failed to update object")
}
