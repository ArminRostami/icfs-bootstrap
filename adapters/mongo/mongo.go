// Package db includes mongoDB implementation of application interfaces
package db

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DBNAME = "icfs"

type mongoDB struct {
	db *mongo.Database
}

func NewMongo(conURI string) (*mongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conURI))
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to mongoDB instance")
	}
	return &mongoDB{db: client.Database(DBNAME)}, nil
}
