package app

import (
	"fmt"
	"icfs_pg/domain"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// TODO: consider removing searchcontent
type ContentStore interface {
	AddContent(c *domain.Content) error
	DeleteContent(id string) error
	GetContent(id string) (*domain.Content, error)
	AddDownload(uid, id string) error
	UpdateContent(id string, updates map[string]interface{}) error
	SearchContent(keys, values []string) (*[]domain.Content, error)
	TextSearch(term string) (*[]domain.Content, error)
	GetAll() (*[]domain.Content, error)
	IncrementDownloads(id string) error
	RateContent(rating float32, uid, cid string) error
	AddComment(uid, id, comment string) error
	GetComments(id string) (*[]domain.Comment, error)
	GetUserContents(uid string) (*[]domain.Content, error)
}

type ContentService struct {
	CST ContentStore
	UST UserStore
}

func (s *ContentService) RegisterContent(c *domain.Content) (string, *Error) {
	c.ID = uuid.New().String()
	c.Downloads = 0
	c.UploadedAt = time.Now()
	c.LastModified = c.UploadedAt
	c.Description = fmt.Sprintf("%.200s", c.Description)
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

// TODO: add transaction support:
func (s *ContentService) GetContentWithID(uid, id string) (*domain.Content, error) {
	downloader, err := s.UST.GetUserWithID(uid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user info")
	}

	content, err := s.CST.GetContent(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get content info")
	}

	if downloader.ID == content.UploaderID {
		return nil, errors.New("the uploader cannot download their own file")
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

	err = s.CST.IncrementDownloads(content.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to increment downloads")
	}

	err = s.CST.AddDownload(downloader.ID, content.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add to downloads")
	}

	return content, err
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

// TODO: consider removing update functionality for contents
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

func (s *ContentService) RateContent(rating float32, uid, cid string) error {
	return s.CST.RateContent(rating, uid, cid)
}

func (s *ContentService) SearchContent(search map[string]string) (*[]domain.Content, error) {
	validKeys := map[string]struct{}{"name": {}, "description": {}, "filename": {}, "extension": {}}
	for key := range search {
		if _, exists := validKeys[key]; !exists {
			delete(search, key)
		}
	}

	results, err := s.CST.SearchContent(getSlicesFromMap(search))
	if err != nil {
		return nil, errors.Wrap(err, "failed to search")
	}
	for i := range *results {
		(*results)[i].CID = ""
	}
	return results, nil
}

func getSlicesFromMap(m map[string]string) ([]string, []string) {
	keys := make([]string, len(m))
	values := make([]string, len(m))

	count := 0
	for k, v := range m {
		keys[count] = k
		values[count] = v
		count++
	}
	return keys, values
}

func (s *ContentService) TextSearch(term string) (*[]domain.Content, error) {
	return s.CST.TextSearch(term)
}

func (s *ContentService) GetAll() (*[]domain.Content, error) {
	return s.CST.GetAll()
}

func (s *ContentService) AddComment(uid, id, comment string) error {
	return s.CST.AddComment(uid, id, comment)
}

func (s *ContentService) GetComments(id string) (*[]domain.Comment, error) {
	return s.CST.GetComments(id)
}

func (s *ContentService) GetUserContents(uid string) (*[]domain.Content, error) {
	return s.CST.GetUserContents(uid)
}
