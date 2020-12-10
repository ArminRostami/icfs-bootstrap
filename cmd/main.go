package main

import (
	db "icfs_cr/adapters/cockroach"
	http "icfs_cr/adapters/http"
	app "icfs_cr/application"
	"log"

	"github.com/pkg/errors"
)

const conStr = "postgres://root:admin@127.0.0.1:34615?sslmode=require"

func run() error {
	crdb, err := db.New(conStr)
	if err != nil {
		return errors.Wrap(err, "failed to create mongo instance")
	}
	userStore := &db.UserStore{CR: crdb}
	userService := &app.UserService{UST: userStore}
	handler := http.Handler{USV: userService}
	return handler.Serve()
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}
