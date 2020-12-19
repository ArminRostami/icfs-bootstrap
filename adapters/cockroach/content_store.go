package crdb

import (
	"icfs_cr/domain"

	"github.com/pkg/errors"
)

type ContentStore struct {
	CR *CRDB
}

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

func (cs *ContentStore) GetCid(id string) (string, error) {
	var cid string
	err := cs.CR.db.Get(&cid, `SELECT cid FROM contents WHERE id=$1`, id)
	if err != nil {
		return "", errors.Wrap(err, "failed to get id")
	}
	return cid, nil
}
