package db

import (
	"icfs_mongo/domain"
	"time"

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

func (us *UserStore) DeleteUser(id string) error {
	return us.MDB.DeleteOne(UsersColl, bson.M{"_id": id})
}

func (us *UserStore) UpdateUser(id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return us.MDB.UpdateOne(UsersColl, bson.M{"_id": id}, bson.M{"$set": updates})
}

func (us *UserStore) GetUserWithID(id string) (*domain.User, error) {
	res := us.MDB.FindOne(UsersColl, bson.M{"_id": id})
	var u domain.User
	err := res.Decode(&u)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode result into user")
	}
	return &u, nil
}
