package adapters

import "mime/multipart"

type FileHeaderAdapter struct {
	Original *multipart.FileHeader
}

func (fha *FileHeaderAdapter) Open() (multipart.File, error) {
	return fha.Original.Open()
}

func (fha *FileHeaderAdapter) Filename() string {
	return fha.Original.Filename
}

func NewFileHeaderAdapter(fh *multipart.FileHeader) *FileHeaderAdapter {
	return &FileHeaderAdapter{Original: fh}
}
