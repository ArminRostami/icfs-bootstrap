package app

import (
	"icfs_mongo/domain"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type ContentStore interface {
	AddContent(c *domain.Content) (string, error)
}

type ContentService struct {
	CST ContentStore
}

func (cs *ContentService) AddContent(c *domain.Content) (string, *Error) {
	c.UploadDate = time.Now()
	id, err := cs.CST.AddContent(c)
	if err != nil {
		return "", &Error{http.StatusInternalServerError, errors.Wrap(err, "failed to add content")}
	}
	return id, nil
}
