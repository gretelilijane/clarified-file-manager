package types

import (
	"time"
)

type File struct {
	ID         int32     `json:"id"`
	UserID     int32     `json:"user_id"`
	Name       string    `json:"name"`
	MimeType   string    `json:"mime_type"`
	UploadedAt time.Time `json:"uploaded_at"`
	Content    []byte    `json:"content"`
	Size       int64     `json:"size"`
}
