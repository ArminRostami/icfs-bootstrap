package db

import (
	"icfs_mongo/domain"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

const UsersColl = "users"

type UserStore struct {
	MDB *mongoDB
}

func (us *UserStore) InsertUser(u *domain.User) (string, error) {
	return us.MDB.InsertOne(UsersColl, u)
}

func (us *UserStore) GetUserWithName(username string) (*domain.User, error) {
	res := us.MDB.FindOne(UsersColl, bson.M{"username": username})
	var u domain.User
	err := res.Decode(&u)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode result into user")
	}
	return &u, nil
}
