package adapters

import "mime/multipart"

type FileHeaderAdapter struct {
	Original *multipart.FileHeader
}

func (fha *FileHeaderAdapter) Open() (multipart.File, error) {
	return fha.Original.Open()
}

func NewFileHeaderAdapter(fh *multipart.FileHeader) *FileHeaderAdapter {
	return &FileHeaderAdapter{Original: fh}
}
