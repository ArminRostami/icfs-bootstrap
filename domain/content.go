package domain

type Content struct {
	CID         string  `json:"cid" db:"cid"`
	Name        string  `json:"name" db:"name"`
	Description string  `json:"description" db:"description"`
	FileName    string  `json:"filename" db:"filename"`
	Extension   string  `json:"extension" db:"extension"`
	Category    string  `json:"category" db:"category"`
	UploaderID  string  `json:"uploader_id" db:"uploader_id"`
	Downloads   int     `json:"downloads" db:"downloads"`
	Size        float32 `json:"size" db:"size"`
}
