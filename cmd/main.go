package main

import (
	http "icfs_mongo/adapters/http"
	db "icfs_mongo/adapters/mongo"
	app "icfs_mongo/application"
	"log"

	"github.com/pkg/errors"
)

const ConStr = "mongodb://root:example@localhost:27017"

func run() error {
	mdb, err := db.NewMongo(ConStr)
	if err != nil {
		return errors.Wrap(err, "failed to create mongo instance")
	}
	userStore := &db.UserStore{MDB: mdb}
	userService := &app.UserService{UST: userStore}
	handler := http.Handler{USV: userService}
	return handler.Serve()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
