package crdb

import (
	"icfs_cr/domain"

	"github.com/pkg/errors"
)

type UserStore struct {
	CR *CRDB
}

func (us *UserStore) InsertUser(user *domain.User) (string, error) {
	query := `
	INSERT INTO users(id, username, password, email, bio, credit, created_at, updated_at)
	VALUES (:id, :username, :password, :email, :bio, :credit, :created_at, :updated_at);`
	rows, err := us.CR.NamedExec(query, user)
	if err != nil || rows <= 0 {
		return "", errors.Wrap(err, "failed to insert user")
	}
	return user.ID, nil
}

func (us *UserStore) GetUserWithName(username string) (*domain.User, error) {
	var user domain.User
	err := us.CR.db.Get(&user, `SELECT * FROM users WHERE username=$1;`, username)
	return &user, errors.Wrap(err, "failed to get user with name")
}

func (us *UserStore) GetUserWithID(id string) (*domain.User, error) {
	var user domain.User
	err := us.CR.db.Get(&user, `SELECT * FROM users WHERE id=$1;`, id)
	return &user, errors.Wrap(err, "failed to get user with id")
}

func (us *UserStore) DeleteUser(id string) error {
	rows, err := us.CR.Exec(`DELETE FROM users WHERE id=$1;`, id)
	if err != nil || rows < 1 {
		return errors.Wrap(err, "failed to delete user")
	}
	return nil
}

func (us *UserStore) UpdateUser(id string, updates map[string]interface{}) error {
	panic("not implemented") // TODO: Implement
}

func (us *UserStore) SearchInBio(term string) (*[]domain.User, error) {
	panic("not implemented") // TODO: Implement
}
