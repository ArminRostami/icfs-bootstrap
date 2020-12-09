// Package crdb includes cockroachDB implementation of application interfaces
package crdb

import (
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type CRDB struct {
	db *sqlx.DB
}

func New(conStr string) (*CRDB, error) {
	dbx, err := sqlx.Connect("pgx", conStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect ot db")
	}
	_, err = dbx.Exec(schema)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute schema")
	}
	return &CRDB{db: dbx}, nil
}

func (c *CRDB) NamedExec(query string, arg interface{}) (int64, error) {
	res, err := c.db.NamedExec(query, arg)
	if err != nil {
		return -1, errors.Wrap(err, "failed to execute named query")
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return -1, errors.Wrap(err, "failed to get affected rows")
	}
	return rows, nil
}

func (c *CRDB) Exec(query string, args ...interface{}) (int64, error) {
	res, err := c.db.Exec(query, args...)
	if err != nil {
		return -1, errors.Wrap(err, "failed to execute query")
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return -1, errors.Wrap(err, "failed to get affected rows")
	}
	return rows, nil
}
