package postgres

import (
	"fmt"
	"icfs_pg/domain"
	"time"

	"github.com/pkg/errors"
)

type ContentStore struct {
	DB *PGSQL
}

func (cs *ContentStore) AddContent(c *domain.Content) error {
	rows, err := cs.DB.NamedExec(`
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

func (cs *ContentStore) GetContent(id string) (*domain.Content, error) {
	var c domain.Content
	err := cs.DB.db.Get(&c, `
	SELECT c.id, c.cid, c.uploader_id, c.name, c.extension, c.description, 
	c.size, c.downloads, c.uploaded_at, c.last_modified, c.rating, f.file_type
	FROM ftypes f left join contents c on f.id = c.type_id 
	WHERE c.id = $1`, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get id")
	}
	return &c, nil
}

func (cs *ContentStore) AddDownload(uid, id string) error {
	rows, err := cs.DB.Exec(`INSERT INTO downloads(user_id, content_id) VALUES($1, $2)`, uid, id)
	if err != nil {
		return errors.Wrap(err, "failed to add download")
	}
	if rows < 1 {
		return errors.New("operation complete but no row was affected")
	}
	return nil
}

func (cs *ContentStore) DeleteContent(id string) error {
	rows, err := cs.DB.Exec(`DELETE FROM contents WHERE id=$1`, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete content")
	}
	if rows < 1 {
		return errors.New("operation complete but no row was affected")
	}
	return nil
}

func (cs *ContentStore) UpdateContent(id string, updates map[string]interface{}) error {
	for key, val := range updates {
		q := fmt.Sprintf(`UPDATE contents SET %s = $1 WHERE id = $2;`, key)
		rows, err := cs.DB.Exec(q, val, id)
		if err != nil {
			return errors.Wrap(err, "failed to update content")
		}
		if rows < 1 {
			return errors.New("operation complete but no row was affected")
		}
	}
	return nil
}

func (cs *ContentStore) SearchContent(keys, values []string) (*[]domain.Content, error) {
	var results []domain.Content
	q := fmt.Sprintf(`
	SELECT c.id, c.cid, c.uploader_id, c.name, c.extension, c.description, 
	c.size, c.downloads, c.uploaded_at, c.last_modified, c.rating, f.file_type
	FROM ftypes f left join contents c on f.id = c.type_id 
	WHERE %s ILIKE $1`, keys[0])
	for i := 1; i < len(keys); i++ {
		q += fmt.Sprintf(` AND %s ILIKE $%d`, keys[i], i+1)
	}

	err := cs.DB.db.Select(&results, q, getInterfaceSlice(values)...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get results")
	}
	return &results, nil
}

func getInterfaceSlice(strs []string) []interface{} {
	r := make([]interface{}, len(strs))
	for idx := range strs {
		r[idx] = "%" + strs[idx] + "%"
	}
	return r
}

func (cs *ContentStore) IncrementDownloads(id string) error {
	rows, err := cs.DB.Exec(`UPDATE contents SET downloads = downloads + 1 WHERE id=$1`, id)
	if err != nil {
		return errors.Wrap(err, "failed to update content")
	}
	if rows < 1 {
		return errors.New("operation complete but no row was affected")
	}
	return nil
}

func (cs *ContentStore) RateContent(rating float32, uid, cid string) error {
	rows, err := cs.DB.Exec(`
	UPDATE downloads set rating=$1
	where downloads.user_id=$2
	and downloads.content_id=$3;
	`, rating, uid, cid)
	if err != nil {
		return errors.Wrap(err, "failed to update content")
	}
	if rows < 1 {
		return errors.New("operation complete but no row was affected")
	}
	return nil
}

func (cs *ContentStore) TextSearch(term string) (*[]domain.Content, error) {
	var results []domain.Content
	q := `
	SELECT c.id, c.uploader_id, c.name, c.extension, c.description, 
	c.size, c.downloads, c.uploaded_at, c.rating, f.file_type
	FROM ftypes f left join contents c on f.id = c.type_id, websearch_to_tsquery('english', $1) query
	WHERE query @@ tsv
	ORDER BY ts_rank_cd(tsv, query) DESC;`
	err := cs.DB.db.Select(&results, q, term)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get results")
	}
	return &results, nil
}

func (cs *ContentStore) GetAll() (*[]domain.Content, error) {
	var results []domain.Content
	q := `
	SELECT c.id, c.uploader_id, c.name, c.extension, c.description, c.size, 
	c.downloads, c.uploaded_at, c.last_modified, c.rating, f.file_type
	FROM contents c join ftypes f on f.id = c.type_id;
	`
	err := cs.DB.db.Select(&results, q)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get results")
	}
	return &results, nil
}

func (cs *ContentStore) AddComment(uid, id, comment string) error {
	rows, err := cs.DB.Exec(`UPDATE downloads SET comment_text=$1,comment_time=$2 
	WHERE user_id=$3 and content_id=$4`, comment, time.Now(), uid, id)
	if err != nil {
		return errors.Wrap(err, "failed to add comment")
	}
	if rows < 1 {
		return errors.New("operation complete but no row was affected")
	}
	return nil
}

func (cs *ContentStore) GetComments(id string) (*[]domain.Comment, error) {
	var comments []domain.Comment
	q := `SELECT d.comment_text, d.rating, d.comment_time, u.username
	FROM (SELECT * from downloads WHERE content_id=$1) d left join users u on d.user_id=u.id;`
	err := cs.DB.db.Select(&comments, q, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get comments")
	}
	return &comments, nil
}
