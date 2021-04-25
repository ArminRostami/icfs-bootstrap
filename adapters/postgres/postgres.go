// Package postgres includes implementation of application interfaces
package postgres

import (
	"fmt"
	"os"
	"strings"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type PGSQL struct {
	db *sqlx.DB
}

func New(user, password, host string) (*PGSQL, error) {
	dbx, err := sqlx.Connect("pgx", getConStr(user, password, host))
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to db")
	}
	schemaBytes, err := os.ReadFile(getSchemaFile("../adapters/postgres/schema.sql"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to open schema file")
	}
	_, err = dbx.Exec(string(schemaBytes))
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute schema")
	}
	return &PGSQL{db: dbx}, nil
}

func getConStr(user, password, host string) string {
	if dockerEnabled() {
		host = "pgsql:5432"
	}
	return fmt.Sprintf("postgres://%s:%s@%s", user, password, host)
}

func getSchemaFile(fileAddr string) string {
	if dockerEnabled() {
		return "./schema.sql"
	}
	return fileAddr
}

func dockerEnabled() bool {
	val, exists := os.LookupEnv("DOCKER_ENABLED")
	if !exists {
		return false
	}
	if strings.EqualFold(val, "1") {
		return true
	}
	return false
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
