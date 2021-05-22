// Package postgres includes implementation of application interfaces
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"icfs_pg/env"
	"os"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type PGSQL struct {
	db *sqlx.DB
}

func New(host string, port int, user, password string) (*PGSQL, error) {
	dbx, err := sqlx.Connect("pgx", getConStr(host, port, user, password))
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

func txFromCtx(ctx context.Context) (*sqlx.Tx, error) {
	tx, ok := ctx.Value(txKey).(*sqlx.Tx)
	if !ok {
		return nil, errors.New("ctx does not include tx")
	}
	return tx, nil
}

func (pg *PGSQL) CtxWithTx() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	tx := pg.db.MustBeginTx(ctx, &sql.TxOptions{})
	return context.WithValue(ctx, txKey, tx), cancel
}

func (pg *PGSQL) TxCommit(ctx context.Context) error {
	tx, err := txFromCtx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get tx from ctx")
	}
	return tx.Commit()
}

func getConStr(host string, port int, user, password string) string {
	if env.DockerEnabled() {
		host = "pgsql"
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%d", user, password, host, port)
}

func getSchemaFile(fileAddr string) string {
	if env.DockerEnabled() {
		return "./schema.sql"
	}
	return fileAddr
}

func NamedExec(tx *sqlx.Tx, query string, arg interface{}) (int64, error) {
	res, err := tx.NamedExec(query, arg)
	if err != nil {
		return -1, errors.Wrap(err, "failed to execute named query")
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return -1, errors.Wrap(err, "failed to get affected rows")
	}
	return rows, nil
}

func Exec(tx *sqlx.Tx, query string, args ...interface{}) (int64, error) {
	res, err := tx.Exec(query, args...)
	if err != nil {
		return -1, errors.Wrap(err, "failed to execute query")
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return -1, errors.Wrap(err, "failed to get affected rows")
	}
	return rows, nil
}
