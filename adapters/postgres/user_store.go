package postgres

import (
	"fmt"
	"icfs_pg/domain"

	"github.com/pkg/errors"
)

type UserStore struct {
	DB *PGSQL
}

const usersTable = "users"

func (us *UserStore) InsertUser(user *domain.User) (string, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s(id, username, password, email, credit, created_at, updated_at)
	VALUES (:id, :username, :password, :email, :credit, :created_at, :updated_at);`, usersTable)
	rows, err := us.DB.NamedExec(query, user)
	if err != nil || rows <= 0 {
		return "", errors.Wrap(err, "failed to insert user")
	}
	return user.ID, nil
}

func (us *UserStore) GetUserWithName(username string) (*domain.User, error) {
	var user domain.User
	query := fmt.Sprintf(`SELECT * FROM %s WHERE username=$1;`, usersTable)
	err := us.DB.db.Get(&user, query, username)
	return &user, errors.Wrap(err, "failed to get user with name")
}

func (us *UserStore) GetUserWithID(id string) (*domain.User, error) {
	var user domain.User
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id=$1;`, usersTable)
	err := us.DB.db.Get(&user, query, id)
	return &user, errors.Wrap(err, "failed to get user with id")
}

func (us *UserStore) DeleteUser(id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1;`, usersTable)
	rows, err := us.DB.Exec(query, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete user")
	}
	if rows < 1 {
		return errors.New("operation complete but no row was affected")
	}
	return nil
}

func (us *UserStore) UpdateUser(id string, updates map[string]interface{}) error {
	for key, val := range updates {
		q := fmt.Sprintf(`UPDATE %s SET %s = $1 WHERE id = $2;`, usersTable, key)
		rows, err := us.DB.Exec(q, val, id)
		if err != nil {
			return errors.Wrap(err, "failed to update user")
		}
		if rows < 1 {
			return errors.New("operation complete but no row was affected")
		}
	}
	return nil
}

func (us *UserStore) ModifyCredit(uid string, value int) error {
	q := fmt.Sprintf(`UPDATE %s SET credit = credit + $1 WHERE id=$2`, usersTable)
	rows, err := us.DB.Exec(q, value, uid)
	if err != nil {
		return errors.Wrap(err, "failed to modify credit")
	}
	if rows < 1 {
		return errors.New("no modification done")
	}
	return nil
}
