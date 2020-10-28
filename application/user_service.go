// Package app contains the application logic
package app

import (
	"icfs_mongo/domain"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const SigningKey = "VhFJdNDsE9vheq6wTEFga7WhuR4TJ1E8JTPNFaH3e_o"

type UserStore interface {
	InsertUser(user *domain.User) (string, error)
	GetUserWithName(username string) (*domain.User, error)
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"username": username, "id": user.ID})
	tokenStr, err := token.SignedString([]byte(SigningKey))
	if err != nil {
		return "", &Error{http.StatusInternalServerError, errors.Wrap(err, "failed to sign jwt")}
	}

	return tokenStr, nil
}

func (s *UserService) ValidateAuth(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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
	return token.Claims, nil
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
