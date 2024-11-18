package types

import (
	"time"
)

type File struct {
	ID         int32
	UserID     int32
	Name       string
	MimeType   string
	UploadedAt time.Time
	Content    []byte
	Size       int64
}

type FileSortableColumn string

const (
	FileSortByName         FileSortableColumn = "name"
	FileSortByMimeType     FileSortableColumn = "mime_type"
	FileSortableSize       FileSortableColumn = "size"
	FileSortableUploadedAt FileSortableColumn = "uploaded_at"
)

func FileSortableColumnFromString(value string) FileSortableColumn {
	column := FileSortableColumn(value)

	switch column {
	case FileSortByName, FileSortByMimeType, FileSortableSize, FileSortableUploadedAt:
		return column
	default:
		return FileSortableUploadedAt
	}
}
