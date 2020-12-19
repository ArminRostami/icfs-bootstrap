package app

import (
	"fmt"
	"icfs_cr/domain"
	"mime"
	"net/http"

	"github.com/google/uuid"
)

type ContentStore interface {
	AddContent(c *domain.Content) error
	GetCid(id string) (string, error)
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

func (s *ContentService) GetContentWithID(id string) (string, error) {
	// check if user has required credit
	// add to uploader and subtract from downloader
	return s.CST.GetCid(id)
}

// add update function

// add delete function: subtract credit from uploader

// add content discovery functions
