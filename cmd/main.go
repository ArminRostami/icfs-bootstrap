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

const localhost = "127.0.0.1"

func run() error {
	pgsql, err := db.New(localhost, 5432, "postgres", "example")
	if err != nil {
		return errors.Wrap(err, "failed to create postgresql instance")
	}

	rds, err := redis.New(localhost, 6379, "")
	if err != nil {
		return errors.Wrap(err, "failed to create reidis instance")
	}

	cancel, service, err := ipfs.NewService()
	defer cancel()
	if err != nil {
		return errors.Wrap(err, "failed to create ipfs service")
	}

	us := &db.UserStore{DB: pgsql}
	cs := &db.ContentStore{DB: pgsql}

	contentService := &app.ContentService{ContentStore: cs, UserStore: us, ContextProvider: pgsql}
	userService := &app.UserService{UserStore: us, SessionStore: rds, ContextProvider: pgsql}

	handler := http.Handler{US: userService, CS: contentService, IS: service}

	return handler.Serve()
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}
