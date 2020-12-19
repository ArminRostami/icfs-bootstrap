package crdb

import (
	"fmt"
	"icfs_cr/domain"

	"github.com/pkg/errors"
)

type ContentStore struct {
	CR *CRDB
}

const contentsTable = "contents"

func (cs *ContentStore) AddContent(c *domain.Content) error {
	rows, err := cs.CR.NamedExec(`
	INSERT INTO contents(id,cid,name,description,filename,extension,category,uploader_id,size,downloads) 
	VALUES(:id,:cid,:name,:description,:filename,:extension,:category,:uploader_id,:size,:downloads) `, c)
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
	err := cs.CR.db.Get(&c, `SELECT * FROM contents WHERE id=$1`, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get id")
	}
	return &c, nil
}

func (cs *ContentStore) DeleteContent(id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1;`, contentsTable)
	rows, err := cs.CR.Exec(query, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete user")
	}
	if rows < 1 {
		return errors.New("operation complete but no row was affected")
	}
	return nil
}
