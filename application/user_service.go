// Package app contains the application logic
package app

import (
	"icfs_mongo/domain"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	RegisterUser(user *domain.User) (string, error)
}

type UserService struct {
	UST UserStore
}

func (s *UserService) RegisterUser(user *domain.User) (string, *Error) {
	hash, err := hashPassword(user.Password)
	if err != nil {
		return "", &Error{http.StatusInternalServerError, err}
	}
	user.Password = hash
	user.Credit = 0
	id, err := s.UST.RegisterUser(user)
	if err != nil {
		return "", &Error{http.StatusInternalServerError, errors.Wrap(err, "failed to register user")}
	}
	return id, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", errors.Wrap(err, "failed to hash password")
	}
	return string(hash), nil
}

// func checkPassword(input, hash string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input))
// 	return err == nil
// }
