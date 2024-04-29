package helpers

import (
	"io"
	"mime/multipart"
	"net/http"
	"vue-converter-backend/interfaces"
)

type FileHeaderAdapter struct {
	FileHeader *multipart.FileHeader
}

func (f *FileHeaderAdapter) Open() (multipart.File, error) {
	return f.FileHeader.Open()
}

type GetTextContentFromFileInterface interface {
	GetTextContentFromFiles(files []interfaces.FileHeader, w http.ResponseWriter) ([]string, error)
}

func GetTextContentFromFiles(files []interfaces.FileHeader, w http.ResponseWriter) ([]string, error) {
	var fileContents []string

	for _, fileHeader := range files {
		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}

		// Read the file content
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return nil, err

		}
		file.Close()
		fileContents = append(fileContents, string(fileBytes))
	}

	return fileContents, nil
}
