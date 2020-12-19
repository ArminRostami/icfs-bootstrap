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
	UpdateContent(id string, updates map[string]interface{}) error
}

type ContentService struct {
	CST ContentStore
	UST UserStore
}

func (s *ContentService) RegisterContent(c *domain.Content) (string, *Error) {
	c.ID = uuid.New().String()
	c.Downloads = 0
	c.Category = mime.TypeByExtension(fmt.Sprintf(".%s", c.Extension))
	err := s.CST.AddContent(c)
	if err != nil {
		return "", &Error{Status: http.StatusInternalServerError, Err: err}
	}
	err = s.UST.ModifyCredit(c.UploaderID, int(c.Size))
	if err != nil {
		return "", &Error{Status: http.StatusInternalServerError, Err: err}
	}
	return c.ID, nil
}

func (s *ContentService) GetContentWithID(uid, id string) (*domain.Content, error) {
	downloader, err := s.UST.GetUserWithID(uid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user info")
	}

	content, err := s.CST.GetContent(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get content info")
	}

	if int(content.Size) > downloader.Credit {
		return nil, errors.Wrap(err, "user does not have enough credit")
	}

	err = s.UST.ModifyCredit(content.UploaderID, int(content.Size))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add credit to uploader")
	}

	err = s.UST.ModifyCredit(uid, -int(content.Size))
	if err != nil {
		return nil, errors.Wrap(err, "failed to subtract credit from downloader")
	}

	return s.CST.GetContent(id)
}

func (s *ContentService) DeleteContent(uid, id string) error {
	c, err := s.CST.GetContent(id)
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

func (s *ContentService) UpdateContent(uid string, updates map[string]interface{}) error {
	id, exists := updates["id"]
	if !exists {
		return errors.New("updates does not include id for content")
	}

	idStr := fmt.Sprint(id)

	c, err := s.CST.GetContent(idStr)
	if err != nil {
		return errors.Wrap(err, "failed to get content")
	}

	if uid != c.UploaderID {
		return errors.New("only the uploader can modify the content")
	}

	validKeys := map[string]struct{}{"name": {}, "description": {}}
	for key := range updates {
		if _, exists := validKeys[key]; !exists {
			delete(updates, key)
		}
	}

	return s.CST.UpdateContent(idStr, updates)

}

// add update function
// add content discovery functions
