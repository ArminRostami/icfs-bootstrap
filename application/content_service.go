package app

import (
	"context"
	"fmt"
	"icfs-boot/domain"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ContentStore interface {
	AddContent(ctx context.Context, c *domain.Content) error
	DeleteContent(ctx context.Context, id string) error
	GetContent(ctx context.Context, id string) (*domain.Content, error)
	AddDownload(ctx context.Context, uid, id string) error
	UpdateContent(ctx context.Context, id string, updates map[string]interface{}) error
	TextSearch(ctx context.Context, term string) (*[]domain.Content, error)
	GetAll(ctx context.Context) (*[]domain.Content, error)
	IncrementDownloads(ctx context.Context, id string) error
	RateContent(ctx context.Context, rating float32, uid, cid string) error
	AddComment(ctx context.Context, uid, id, comment string) error
	GetComments(ctx context.Context, id string) (*[]domain.Comment, error)
	GetUserContents(ctx context.Context, uid string) (*[]domain.Content, error)
}

type ContentService struct {
	ContentStore
	UserStore
	ContextProvider
}

func (s *ContentService) RegisterContent(c *domain.Content) (string, *Error) {
	c.ID = uuid.New().String()
	c.Downloads = 0
	c.UploadedAt = time.Now()
	c.LastModified = c.UploadedAt
	c.Description = fmt.Sprintf("%.200s", c.Description)

	ctx, cancel := s.CtxWithTx()
	defer cancel()

	err := s.AddContent(ctx, c)
	if err != nil {
		return "", &Error{Status: http.StatusInternalServerError, Err: err}
	}

	err = s.ModifyCredit(ctx, c.UploaderID, int(c.Size))
	if err != nil {
		return "", &Error{Status: http.StatusInternalServerError, Err: err}
	}

	if err = s.TxCommit(ctx); err != nil {
		return "", &Error{Status: http.StatusInternalServerError, Err: err}
	}

	return c.ID, nil
}

func (s *ContentService) GetContentWithID(uid, id string) (*domain.Content, error) {
	ctx, cancel := s.CtxWithTx()
	defer cancel()

	downloader, err := s.GetUserWithID(ctx, uid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user info")
	}

	content, err := s.GetContent(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get content info")
	}

	if downloader.ID == content.UploaderID {
		return nil, errors.New("the uploader cannot download their own file")
	}

	if int(content.Size) > downloader.Credit {
		return nil, errors.Wrap(err, "user does not have enough credit")
	}

	err = s.ModifyCredit(ctx, content.UploaderID, int(content.Size))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add credit to uploader")
	}

	err = s.ModifyCredit(ctx, uid, -int(content.Size))
	if err != nil {
		return nil, errors.Wrap(err, "failed to subtract credit from downloader")
	}

	err = s.IncrementDownloads(ctx, content.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to increment downloads")
	}

	err = s.AddDownload(ctx, downloader.ID, content.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add to downloads")
	}

	if err = s.TxCommit(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to commit tx")
	}

	return content, err
}

func (s *ContentService) DeleteContent(uid, id string) error {
	ctx, cancel := s.CtxWithTx()
	defer cancel()

	c, err := s.GetContent(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to get content id")
	}

	if uid != c.UploaderID {
		return errors.New("failed to delete: only the uploader can delete file")
	}

	err = s.ModifyCredit(ctx, uid, -int(c.Size))
	if err != nil {
		return errors.Wrap(err, "failed to decrease credit")
	}

	err = s.ContentStore.DeleteContent(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete content")
	}

	if err = s.TxCommit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit tx")
	}

	return nil
}

// TODO: consider removing update functionality for contents
func (s *ContentService) UpdateContent(uid string, updates map[string]interface{}) error {
	id, exists := updates["id"]
	if !exists {
		return errors.New("updates does not include id for content")
	}

	idStr := fmt.Sprint(id)

	ctx, cancel := s.CtxWithTx()
	defer cancel()

	c, err := s.GetContent(ctx, idStr)
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

	err = s.ContentStore.UpdateContent(ctx, idStr, updates)
	if err != nil {
		return errors.Wrap(err, "failed to update content")
	}

	if err = s.TxCommit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit tx")
	}

	return nil

}

func (s *ContentService) RateContent(rating float32, uid, cid string) error {
	ctx, cancel := s.CtxWithTx()
	defer cancel()

	err := s.ContentStore.RateContent(ctx, rating, uid, cid)
	if err != nil {
		return errors.Wrapf(err, "failed to rate content")
	}

	if err = s.TxCommit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit tx")
	}

	return nil
}

func (s *ContentService) TextSearch(term string) (*[]domain.Content, error) {
	ctx, cancel := s.CtxWithTx()
	defer cancel()

	content, err := s.ContentStore.TextSearch(ctx, term)
	if err != nil {
		return nil, errors.Wrap(err, "failed to search content")
	}

	if err = s.TxCommit(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to commit tx")
	}
	return content, nil
}

func (s *ContentService) GetAll() (*[]domain.Content, error) {
	ctx, cancel := s.CtxWithTx()
	defer cancel()

	contents, err := s.ContentStore.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user contents")
	}

	if err = s.TxCommit(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to commit tx")
	}

	return contents, nil
}

func (s *ContentService) AddComment(uid, id, comment string) error {
	ctx, cancel := s.CtxWithTx()
	defer cancel()

	err := s.ContentStore.AddComment(ctx, uid, id, comment)
	if err != nil {
		return errors.Wrap(err, "failed to add comment")
	}

	if err = s.TxCommit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit tx")
	}
	return nil
}

func (s *ContentService) GetComments(id string) (*[]domain.Comment, error) {
	ctx, cancel := s.CtxWithTx()
	defer cancel()

	comments, err := s.ContentStore.GetComments(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get comments")
	}

	if err = s.TxCommit(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to commit tx")
	}

	return comments, nil
}

func (s *ContentService) GetUserContents(uid string) (*[]domain.Content, error) {
	ctx, cancel := s.CtxWithTx()
	defer cancel()

	contents, err := s.ContentStore.GetUserContents(ctx, uid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user contents")
	}

	if err = s.TxCommit(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to commit tx")
	}
	return contents, nil
}
