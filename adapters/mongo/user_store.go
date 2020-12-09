package db

import (
	"context"
	"icfs_cr/domain"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const UsersColl = "users"

type UserStore struct {
	MCL  *MongoCol
	once sync.Once
}

func NewUserStore(mdb *mongo.Database) *UserStore {
	return &UserStore{MCL: &MongoCol{Col: mdb.Collection(UsersColl)}}
}

func (us *UserStore) CreateIndexes() {
	idxs := []mongo.IndexModel{
		{Keys: bson.M{"bio": "text"}},
		{Keys: bson.M{"username": 1}, Options: options.Index().SetUnique(true)},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = us.MCL.Col.Indexes().CreateMany(ctx, idxs)
}

func (us *UserStore) InsertUser(u *domain.User) (string, error) {
	id, err := us.MCL.InsertOne(u)
	if err != nil {
		return "", errors.Wrap(err, "failed to insert user")
	}
	us.once.Do(us.CreateIndexes)
	return id, nil
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

func (us *UserStore) SearchInBio(term string) (*[]domain.User, error) {
	cur, err := us.MCL.Find(bson.M{"$text": bson.M{"$search": term}})
	if err != nil {
		return nil, errors.Wrap(err, "no matches")
	}
	var results []domain.User
	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()
	err = cur.All(ctx, &results)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode search results")
	}
	return &results, nil
}
