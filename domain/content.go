package domain

import "time"

type Content struct {
	Title       string    `bson:"title" json:"title" binding:"required"`
	Description string    `bson:"description" json:"description"`
	UploadDate  time.Time `bson:"upload_date" json:"upload_date"`
	Rating      float32   `bson:"rating" json:"rating"`
	Type        string    `bson:"type" json:"type"`
	Comments    []string  `bson:"comments" json:"comments"`
	Cid         string    `bson:"cid" json:"cid"`
	Downloads   int       `bson:"downloads" json:"downloads"`
	Size        float32   `bson:"size" json:"size"`
	Uploader    User      `bson:"uploader" json:"uploader"`
}
