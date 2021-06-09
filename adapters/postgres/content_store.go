package postgres

import (
	"context"
	"fmt"
	"icfs-boot/domain"
	"time"

	"github.com/pkg/errors"
)

type ContentStore struct {
	DB *PGSQL
}

type ctxkey int

var txKey ctxkey = 0

func (cs *ContentStore) AddContent(ctx context.Context, c *domain.Content) error {
	tx, err := txFromCtx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get tx from ctx")
	}

	rows, err := NamedExec(tx, `
	INSERT INTO contents(id,cid,name,description,extension,type_id,uploader_id,size,downloads) 
	VALUES(:id,:cid,:name,:description,:extension,(SELECT id from ftypes where file_type=:file_type),:uploader_id,:size,:downloads) `, c)
	if err != nil {
		return errors.Wrap(err, "failed to add content")
	}
	if rows < 1 {
		return errors.New("content was not added")
	}
	return nil
}

func (cs *ContentStore) GetContent(ctx context.Context, id string) (*domain.Content, error) {
	tx, err := txFromCtx(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tx from ctx")
	}
	var c domain.Content
	err = tx.Get(&c, `
	SELECT c.id, c.cid, c.uploader_id, c.name, c.extension, c.description, 
	c.size, c.downloads, c.uploaded_at, c.last_modified, c.rating, f.file_type
	FROM ftypes f left join contents c on f.id = c.type_id 
	WHERE c.id = $1`, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get id")
	}
	return &c, nil
}

func (cs *ContentStore) AddDownload(ctx context.Context, uid, id string) error {
	tx, err := txFromCtx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get tx from ctx")
	}

	_, err = Exec(tx, `
	INSERT INTO downloads(user_id, content_id) VALUES($1, $2) 
	ON CONFLICT ON CONSTRAINT unique_ratings DO NOTHING`, uid, id)
	if err != nil {
		return errors.Wrap(err, "failed to add download")
	}
	return nil
}

func (cs *ContentStore) DeleteContent(ctx context.Context, id string) error {
	tx, err := txFromCtx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get tx from ctx")
	}

	rows, err := Exec(tx, `DELETE FROM contents WHERE id=$1`, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete content")
	}
	if rows < 1 {
		return errors.New("operation complete but no row was affected")
	}
	return nil
}

func (cs *ContentStore) DeleteDownload(ctx context.Context, uid, id string) error {
	tx, err := txFromCtx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get tx from ctx")
	}

	rows, err := Exec(tx, `DELETE FROM downloads WHERE content_id=$1 and user_id=$2`, id, uid)
	if err != nil {
		return errors.Wrap(err, "failed to delete content")
	}
	if rows < 1 {
		return errors.New("operation complete but no row was affected")
	}
	return nil
}

func (cs *ContentStore) UpdateContent(ctx context.Context, id string, updates map[string]interface{}) error {
	tx, err := txFromCtx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get tx from ctx")
	}
	for key, val := range updates {
		q := fmt.Sprintf(`UPDATE contents SET %s = $1 WHERE id = $2;`, key)
		rows, err := Exec(tx, q, val, id)
		if err != nil {
			return errors.Wrap(err, "failed to update content")
		}
		if rows < 1 {
			return errors.New("operation complete but no row was affected")
		}
	}
	return nil
}

func (cs *ContentStore) IncrementDownloads(ctx context.Context, id string) error {
	tx, err := txFromCtx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get tx from ctx")
	}
	rows, err := Exec(tx, `UPDATE contents SET downloads = downloads + 1 WHERE id=$1`, id)
	if err != nil {
		return errors.Wrap(err, "failed to update content")
	}
	if rows < 1 {
		return errors.New("operation complete but no row was affected")
	}
	return nil
}

func (cs *ContentStore) TextSearch(ctx context.Context, term string) (*[]domain.Content, error) {
	tx, err := txFromCtx(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tx from ctx")
	}
	var results []domain.Content
	q := `
	SELECT c.id, c.uploader_id, c.name, c.extension, c.description, 
	c.size, c.downloads, c.uploaded_at, c.rating, f.file_type
	FROM ftypes f left join contents c on f.id = c.type_id, websearch_to_tsquery('english', $1) query
	WHERE query @@ tsv
	ORDER BY ts_rank_cd(tsv, query) DESC;`
	err = tx.Select(&results, q, term)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get results")
	}
	return &results, nil
}

func (cs *ContentStore) GetAll(ctx context.Context) (*[]domain.Content, error) {
	tx, err := txFromCtx(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tx from ctx")
	}

	var results []domain.Content
	q := `
	SELECT c.id, c.uploader_id, c.name, c.extension, c.description, c.size, 
	c.downloads, c.uploaded_at, c.last_modified, c.rating, f.file_type
	FROM contents c join ftypes f on f.id = c.type_id;
	`
	err = tx.Select(&results, q)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get results")
	}
	return &results, nil
}

func (cs *ContentStore) AddReview(ctx context.Context, uid, id, comment string, rating float32) error {
	tx, err := txFromCtx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get tx from ctx")
	}
	rows, err := Exec(tx, `UPDATE downloads SET comment_text=$1,comment_time=$2,rating=$3
	WHERE user_id=$4 and content_id=$5`, comment, time.Now(), rating, uid, id)
	if err != nil {
		return errors.Wrap(err, "failed to add comment")
	}
	if rows < 1 {
		return errors.New("operation complete but no row was affected")
	}
	return nil

}

func (cs *ContentStore) GetComments(ctx context.Context, id string) (*[]domain.Comment, error) {
	tx, err := txFromCtx(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tx from ctx")
	}

	var comments []domain.Comment
	q := `SELECT d.comment_text, d.rating, d.comment_time, u.username
	FROM (SELECT * from downloads WHERE content_id=$1) d left join users u on d.user_id=u.id;`
	err = tx.Select(&comments, q, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get comments")
	}
	return &comments, nil
}

func (cs *ContentStore) GetUserUploads(ctx context.Context, uid string) (*[]domain.Content, error) {
	tx, err := txFromCtx(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tx from ctx")
	}

	var results []domain.Content
	q := `
	SELECT c.id, c.cid, c.uploader_id, c.name, c.extension, c.description, c.size, 
	c.downloads, c.uploaded_at, c.last_modified, c.rating, f.file_type
	FROM contents c join ftypes f on f.id = c.type_id
	WHERE c.uploader_id = $1;
	`
	err = tx.Select(&results, q, uid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get results")
	}
	return &results, nil
}

func (cs *ContentStore) GetUserDownloads(ctx context.Context, uid string) (*[]domain.Content, error) {
	tx, err := txFromCtx(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tx from ctx")
	}

	var results []domain.Content

	q := `SELECT c.id, c.cid, c.uploader_id, c.name, c.extension, c.description, 
	c.size, c.downloads, c.rating, c.uploaded_at, c.last_modified, f.file_type
	FROM (select content_id from downloads where user_id = $1) as d 
	left join contents c on d.content_id = c.id left join ftypes f on c.type_id = f.id`

	err = tx.Select(&results, q, uid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get results")
	}
	return &results, nil
}
