package app

import (
	"fmt"
	"icfs_cr/domain"
	"mime"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ContentStore interface {
	AddContent(c *domain.Content) error
	DeleteContent(id string) error
	GetContent(id string) (*domain.Content, error)
}

type ContentService struct {
	CST ContentStore
	UST UserStore
}

func (s *ContentService) RegisterContent(c *domain.Content) *Error {
	c.ID = uuid.New().String()
	c.Downloads = 0
	c.Category = mime.TypeByExtension(fmt.Sprintf(".%s", c.Extension))
	err := s.CST.AddContent(c)
	if err != nil {
		return &Error{Status: http.StatusInternalServerError, Err: err}
	}
	err = s.UST.ModifyCredit(c.UploaderID, int(c.Size))
	if err != nil {
		return &Error{Status: http.StatusInternalServerError, Err: err}
	}
	return nil
}

func (s *ContentService) GetContentWithID(id string) (*domain.Content, error) {
	// check if user has required credit
	// add to uploader and subtract from downloader
	return s.CST.GetContent(id)
}

func (s *ContentService) DeleteContent(uid, id string) error {
	c, err := s.GetContentWithID(id)
	if err != nil {
		return errors.Wrap(err, "failed to get content id")
	}
	if uid != c.UploaderID {
		return errors.New("failed to delete: only the uploader can delete file")
	}
	err = s.UST.ModifyCredit(uid, -int(c.Size))
	if err != nil {
		return errors.Wrap(err, "failed to decrease credit")
	}
	return s.CST.DeleteContent(id)
}

// add update function
// add content discovery functions
