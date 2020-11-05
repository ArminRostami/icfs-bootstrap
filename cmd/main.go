package main

import (
	http "icfs_mongo/adapters/http"
	db "icfs_mongo/adapters/mongo"
	app "icfs_mongo/application"
	"log"

	"github.com/pkg/errors"
)

const ConStr = "mongodb://root:example@localhost:27017"
const UsersColl = "users"

func run() error {
	mdb, err := db.NewMongo(ConStr)
	if err != nil {
		return errors.Wrap(err, "failed to create mongo instance")
	}
	userCollection := &db.MongoCol{Col: mdb.Collection(UsersColl)}
	userStore := &db.UserStore{MCL: userCollection}
	userService := &app.UserService{UST: userStore}
	handler := http.Handler{USV: userService}
	return handler.Serve()
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}
