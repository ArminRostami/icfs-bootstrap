package db

import "icfs_mongo/domain"

const ContentsColl = "contents"

type ContentStore struct {
	MDB *mongoDB
}

func (cs *ContentStore) AddContent(c *domain.Content) (string, error) {
	return cs.MDB.InsertOne(ContentsColl, c)
}
