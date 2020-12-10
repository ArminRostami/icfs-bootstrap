package app

import (
	"fmt"
	"icfs_cr/domain"
	"mime"
	"net/http"
)

type ContentStore interface {
	AddContent(c *domain.Content) error
}

type ContentService struct {
	CST ContentStore
	UST UserStore
}

func (s *ContentService) RegisterContent(c *domain.Content) *Error {
	c.Downloads = 0
	c.Category = mime.TypeByExtension(fmt.Sprintf(".%s", c.Extension))
	err := s.CST.AddContent(c)
	if err != nil {
		return &Error{Status: http.StatusInternalServerError, Err: err}
	}
	return nil
}
