package db

import (
	"icfs_mongo/domain"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type UserStore struct {
	MCL *MongoCol
}

func (us *UserStore) InsertUser(u *domain.User) (string, error) {
	return us.MCL.InsertOne(u)
}

func (us *UserStore) GetUserWithName(username string) (*domain.User, error) {
	res := us.MCL.FindOne(bson.M{"username": username})
	var u domain.User
	err := res.Decode(&u)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode result into user")
	}
	return &u, nil
}

func (us *UserStore) DeleteUser(id string) error {
	return us.MCL.DeleteOne(bson.M{"_id": id})
}

func (us *UserStore) UpdateUser(id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return us.MCL.UpdateOne(bson.M{"_id": id}, bson.M{"$set": updates})
}

func (us *UserStore) GetUserWithID(id string) (*domain.User, error) {
	res := us.MCL.FindOne(bson.M{"_id": id})
	var u domain.User
	err := res.Decode(&u)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode result into user")
	}
	return &u, nil
}
