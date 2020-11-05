// Package app contains the application logic
package app

import (
	"fmt"
	"icfs_mongo/domain"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const SigningKey = "VhFJdNDsE9vheq6wTEFga7WhuR4TJ1E8JTPNFaH3e_o"

type UserStore interface {
	InsertUser(user *domain.User) (string, error)
	GetUserWithName(username string) (*domain.User, error)
	GetUserWithID(id string) (*domain.User, error)
	DeleteUser(id string) error
	UpdateUser(id string, updates map[string]interface{}) error
}

type UserService struct {
	UST UserStore
}

type CustomClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

func (s *UserService) RegisterUser(user *domain.User) (string, *Error) {
	hash, err := hashPassword(user.Password)
	if err != nil {
		return "", &Error{http.StatusInternalServerError, err}
	}
	user.Password = hash
	user.Credit = 0
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	id, err := s.UST.InsertUser(user)
	if err != nil {
		return "", &Error{http.StatusInternalServerError, errors.Wrap(err, "failed to register user")}
	}
	return id, nil
}

func (s *UserService) AuthenticateUser(username, password string) (string, *Error) {
	user, err := s.UST.GetUserWithName(username)
	if err != nil {
		return "", &Error{http.StatusUnauthorized, errors.Wrap(err, "failed to get user from db")}
	}

	if match := checkPassword(password, user.Password); !match {
		return "", &Error{http.StatusUnauthorized, errors.New("auth failed")}
	}
	claims := CustomClaims{ID: user.ID}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(SigningKey))
	if err != nil {
		return "", &Error{http.StatusInternalServerError, errors.Wrap(err, "failed to sign jwt")}
	}

	return tokenStr, nil
}

func (s *UserService) ValidateAuth(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SigningKey), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse token")
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (s *UserService) GetUserWithID(id string) (*domain.User, error) {
	u, err := s.UST.GetUserWithID(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user from userstore")
	}
	u.Password = ""
	return u, nil
}

func (s *UserService) DeleteUser(id string) error {
	return s.UST.DeleteUser(id)
}

func (s *UserService) UpdateUser(id string, updates map[string]interface{}) error {
	if pass, exists := updates["password"]; exists {
		hashed, err := hashPassword(fmt.Sprint(pass))
		if err != nil {
			return errors.Wrap(err, "failed to hash password")
		}
		updates["password"] = hashed
	}

	validKeys := map[string]struct{}{"password": {}, "email": {}, "bio": {}}
	for key := range updates {
		if _, exists := validKeys[key]; !exists {
			delete(updates, key)
		}
	}
	return s.UST.UpdateUser(id, updates)
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
