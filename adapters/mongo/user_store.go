package db

import (
	"context"
	"fmt"
	"icfs_mongo/domain"
	"time"

	"github.com/pkg/errors"
)

const UsersColl = "users"

type UserStore struct {
	MDB *mongoDB
}

func (us *UserStore) RegisterUser(u *domain.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := us.MDB.db.Collection(UsersColl).InsertOne(ctx, u)
	return fmt.Sprintf("%v", res.InsertedID), errors.Wrap(err, "failed to insert user into mongoDB")
}
