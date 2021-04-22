package main

import (
	http "icfs_pg/adapters/http"
	db "icfs_pg/adapters/postgres"
	app "icfs_pg/application"
	"icfs_pg/env"
	"log"

	"github.com/pkg/errors"
)

const conStr = "postgres://postgres:example@127.0.0.1:5432"

func run() error {
	pgsql, err := db.New(getConStr())
	if err != nil {
		return errors.Wrap(err, "failed to create postgresql instance")
	}
	userStore := &db.UserStore{DB: pgsql}
	contentStore := &db.ContentStore{DB: pgsql}
	contentService := &app.ContentService{CST: contentStore, UST: userStore}
	userService := &app.UserService{UST: userStore}
	handler := http.Handler{US: userService, CS: contentService}
	return handler.Serve()
}

func getConStr() string {
	if env.DockerEnabled() {
		return "postgres://postgres:example@pgsql:5432"
	}
	return conStr
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}
