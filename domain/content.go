package domain

import "time"

type Content struct {
	ID           string    `json:"id" db:"id"`
	CID          string    `json:"cid" db:"cid"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	Extension    string    `json:"extension" db:"extension"`
	FileType     string    `json:"file_type" db:"file_type"`
	UploaderID   string    `json:"uploader_id" db:"uploader_id"`
	Downloads    int       `json:"downloads" db:"downloads"`
	Rating       float32   `json:"rating" db:"rating"`
	Size         float32   `json:"size" db:"size"`
	UploadedAt   time.Time `json:"uploaded_at" db:"uploaded_at"`
	LastModified time.Time `json:"last_modified" db:"last_modified"`
}

type Comment struct {
	Username string  `json:"username" db:"username"`
	Rating   float32 `json:"rating" db:"rating"`
	CText    string  `json:"comment_text" db:"comment_text"`
	CTime    string  `json:"comment_time" db:"comment_time"`
}
