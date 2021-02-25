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
	Size         float32   `json:"size" db:"size"`
	UploadedAt   time.Time `json:"uploaded_at" db:"uploaded_at"`
	LastModified time.Time `json:"last_modified" db:"last_modified"`
}

// TODO: add rating
