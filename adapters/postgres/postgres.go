// Package crdb includes cockroachDB implementation of application interfaces
package postgres

import (
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type PGSQL struct {
	db *sqlx.DB
}

func New(conStr string) (*PGSQL, error) {
	dbx, err := sqlx.Connect("pgx", conStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to db")
	}
	_, err = dbx.Exec(schema)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute schema")
	}
	return &PGSQL{db: dbx}, nil
}

func (c *PGSQL) NamedExec(query string, arg interface{}) (int64, error) {
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

func (c *PGSQL) Exec(query string, args ...interface{}) (int64, error) {
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
