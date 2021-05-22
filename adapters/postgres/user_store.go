package postgres

import (
	"context"
	"fmt"
	"icfs_pg/domain"

	"github.com/pkg/errors"
)

type UserStore struct {
	DB *PGSQL
}

const usersTable = "users"

func (us *UserStore) InsertUser(ctx context.Context, user *domain.User) (string, error) {
	tx, err := txFromCtx(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to get tx from ctx")
	}

	query := fmt.Sprintf(`
	INSERT INTO %s(id, username, password, email, credit, created_at, updated_at)
	VALUES (:id, :username, :password, :email, :credit, :created_at, :updated_at);`, usersTable)
	rows, err := NamedExec(tx, query, user)
	if err != nil || rows <= 0 {
		return "", errors.Wrap(err, "failed to insert user")
	}
	return user.ID, nil
}

func (us *UserStore) GetUserWithName(ctx context.Context, username string) (*domain.User, error) {
	tx, err := txFromCtx(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tx from ctx")
	}

	var user domain.User
	query := fmt.Sprintf(`SELECT * FROM %s WHERE username=$1;`, usersTable)
	err = tx.Get(&user, query, username)
	return &user, errors.Wrap(err, "failed to get user with name")
}

func (us *UserStore) GetUserWithID(ctx context.Context, id string) (*domain.User, error) {
	tx, err := txFromCtx(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tx from ctx")
	}

	var user domain.User
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id=$1;`, usersTable)
	err = tx.Get(&user, query, id)
	return &user, errors.Wrap(err, "failed to get user with id")
}

func (us *UserStore) DeleteUser(ctx context.Context, id string) error {
	tx, err := txFromCtx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get tx from ctx")
	}

	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1;`, usersTable)
	rows, err := Exec(tx, query, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete user")
	}
	if rows < 1 {
		return errors.New("operation complete but no row was affected")
	}
	return nil
}

func (us *UserStore) UpdateUser(ctx context.Context, id string, updates map[string]interface{}) error {
	tx, err := txFromCtx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get tx from ctx")
	}

	for key, val := range updates {
		q := fmt.Sprintf(`UPDATE %s SET %s = $1 WHERE id = $2;`, usersTable, key)
		rows, err := Exec(tx, q, val, id)
		if err != nil {
			return errors.Wrap(err, "failed to update user")
		}
		if rows < 1 {
			return errors.New("operation complete but no row was affected")
		}
	}
	return nil
}

func (us *UserStore) ModifyCredit(ctx context.Context, uid string, value int) error {
	tx, err := txFromCtx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get tx from ctx")
	}

	q := fmt.Sprintf(`UPDATE %s SET credit = credit + $1 WHERE id=$2`, usersTable)
	rows, err := Exec(tx, q, value, uid)
	if err != nil {
		return errors.Wrap(err, "failed to modify credit")
	}
	if rows < 1 {
		return errors.New("no modification done")
	}
	return nil
}
