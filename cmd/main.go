package main

import (
	http "icfs_pg/adapters/http"
	"icfs_pg/adapters/ipfs"
	db "icfs_pg/adapters/postgres"
	"icfs_pg/adapters/redis"
	app "icfs_pg/application"
	"log"

	"github.com/pkg/errors"
)

func run() error {
	pgsql, err := db.New("postgres", "example", "127.0.0.1:5432")
	if err != nil {
		return errors.Wrap(err, "failed to create postgresql instance")
	}

	rds := redis.New("localhost:6379", "")

	cancel, service, err := ipfs.NewService()
	defer cancel()
	if err != nil {
		return errors.Wrap(err, "failed to create ipfs service")
	}

	userStore := &db.UserStore{DB: pgsql}
	contentStore := &db.ContentStore{DB: pgsql}

	contentService := &app.ContentService{CST: contentStore, UST: userStore}
	userService := &app.UserService{UST: userStore, SST: rds}

	handler := http.Handler{US: userService, CS: contentService, IS: service}

	return handler.Serve()
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}
