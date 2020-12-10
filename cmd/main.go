package main

import (
	db "icfs_cr/adapters/cockroach"
	http "icfs_cr/adapters/http"
	app "icfs_cr/application"
	"log"

	"github.com/pkg/errors"
)

const conStr = "postgres://root:admin@127.0.0.1:44275?sslmode=require"

func run() error {
	crdb, err := db.New(conStr)
	if err != nil {
		return errors.Wrap(err, "failed to create mongo instance")
	}
	userStore := &db.UserStore{CR: crdb}
	contentStore := &db.ContentStore{CR: crdb}
	contentService := &app.ContentService{CST: contentStore}
	userService := &app.UserService{UST: userStore}
	handler := http.Handler{USV: userService, CS: contentService}
	return handler.Serve()
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}
