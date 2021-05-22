// Package app contains the application logic
package app

import (
	"context"
	"fmt"
	"icfs_pg/domain"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const SigningKey = "VhFJdNDsE9vheq6wTEFga7WhuR4TJ1E8JTPNFaH3e_o"

type UserStore interface {
	InsertUser(ctx context.Context, user *domain.User) (string, error)
	GetUserWithName(ctx context.Context, username string) (*domain.User, error)
	GetUserWithID(ctx context.Context, id string) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id string, updates map[string]interface{}) error
	ModifyCredit(ctx context.Context, uid string, value int) error
}

type SessionStore interface {
	Get(key string) (string, error)
	SetEx(key, value string, expiration int64) error
	Del(key string) error
}

type UserService struct {
	UserStore
	SessionStore
	ContextProvider
}

func (s *UserService) RegisterUser(user *domain.User) (string, *Error) {
	user.ID = uuid.New().String()

	hash, err := hashPassword(user.Password)
	if err != nil {
		return "", &Error{http.StatusInternalServerError, err}
	}

	user.Password = hash
	user.Credit = 0
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	ctx, cancel := s.CtxWithTx()
	defer cancel()

	id, err := s.InsertUser(ctx, user)
	if err != nil {
		return "", &Error{http.StatusInternalServerError, errors.Wrap(err, "failed to register user")}
	}

	if err = s.TxCommit(ctx); err != nil {
		return "", &Error{http.StatusInternalServerError, errors.Wrap(err, "failed to commit TX")}
	}

	return id, nil
}

func (s *UserService) AuthenticateUser(username, password string) (*domain.User, string, *Error) {
	ctx, cancel := s.CtxWithTx()
	defer cancel()

	user, err := s.GetUserWithName(ctx, username)
	if err != nil {
		return nil, "", &Error{http.StatusUnauthorized, errors.Wrap(err, "failed to get user from db")}
	}

	if match := checkPassword(password, user.Password); !match {
		return nil, "", &Error{http.StatusUnauthorized, errors.New("auth failed")}
	}

	sessID := uuid.New().String()
	err = s.SetEx(sessID, user.ID, 24*3600)
	if err != nil {
		return nil, "", &Error{http.StatusInternalServerError, errors.Wrap(err, "failed to set sessID")}
	}

	user.Password = ""
	return user, sessID, nil
}

func (s *UserService) ValidateAuth(sessID string) (string, error) {
	return s.Get(sessID)
}

func (s *UserService) GetUserWithID(id string) (*domain.User, error) {
	ctx, cancel := s.CtxWithTx()
	defer cancel()

	u, err := s.UserStore.GetUserWithID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user from userstore")
	}

	u.Password = ""
	return u, nil
}

func (s *UserService) Logout(sessID string) error {
	return s.Del(sessID)
}

func (s *UserService) DeleteUser(id string) error {
	ctx, cancel := s.CtxWithTx()
	defer cancel()

	err := s.UserStore.DeleteUser(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete user")
	}

	if err = s.TxCommit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit TX")
	}

	return nil
}

func (s *UserService) UpdateUser(id string, updates map[string]interface{}) error {
	if pass, exists := updates["password"]; exists {
		hashed, err := hashPassword(fmt.Sprint(pass))
		if err != nil {
			return errors.Wrap(err, "failed to hash password")
		}
		updates["password"] = hashed
	}

	validKeys := map[string]struct{}{"password": {}, "email": {}}
	for key := range updates {
		if _, exists := validKeys[key]; !exists {
			delete(updates, key)
		}
	}
	ctx, cancel := s.CtxWithTx()
	defer cancel()

	err := s.UserStore.UpdateUser(ctx, id, updates)
	if err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	if err = s.TxCommit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit TX")
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", errors.Wrap(err, "failed to hash password")
	}
	return string(hash), nil
}

func checkPassword(input, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(input)) == nil
}
