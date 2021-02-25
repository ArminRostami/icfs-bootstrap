package postgres

import (
	"fmt"
	"icfs_pg/domain"

	"github.com/pkg/errors"
)

type ContentStore struct {
	DB *PGSQL
}

const contentsTable = "contents"

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
	c.size, c.downloads, c.uploaded_at, c.last_modified, f.file_type
	FROM ftypes f left join contents c on f.id = c.type_id 
	WHERE c.id = $1`, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get id")
	}
	return &c, nil
}

func (cs *ContentStore) DeleteContent(id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1;`, contentsTable)
	rows, err := cs.DB.Exec(query, id)
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
		q := fmt.Sprintf(`UPDATE %s SET %s = $1 WHERE id = $2;`, contentsTable, key)
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
	c.size, c.downloads, c.uploaded_at, c.last_modified, f.file_type
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
	q := fmt.Sprintf(`UPDATE %s SET downloads = downloads + 1 WHERE id=$1`, contentsTable)
	rows, err := cs.DB.Exec(q, id)
	if err != nil {
		return errors.Wrap(err, "failed to update content")
	}
	if rows < 1 {
		return errors.New("operation complete but no row was affected")
	}
	return nil
}
