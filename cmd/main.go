package main

import (
	"context"
	http "icfs_pg/adapters/http"
	"icfs_pg/adapters/ipfs"
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

	cancel, err := initIPFS()
	if err != nil {
		return errors.Wrap(err, "failed to init ipfs service")
	}
	defer cancel()

	return handler.Serve()
}

func initIPFS() (context.CancelFunc, error) {
	cancel, service, err := ipfs.NewService()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create ipfs service")
	}
	err = service.Start()
	if err != nil {
		return nil, errors.Wrap(err, "failed to start ipfs service")
	}
	return cancel, nil
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
