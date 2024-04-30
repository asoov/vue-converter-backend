package helpers

import (
	"io"
	"mime/multipart"
	"net/http"
	"vue-converter-backend/interfaces"
	"vue-converter-backend/models"
)

type FileHeaderAdapter struct {
	FileHeader *multipart.FileHeader
}

func (f *FileHeaderAdapter) Open() (multipart.File, error) {
	return f.FileHeader.Open()
}

type GetTextContentFromFileInterface interface {
	GetTextContentFromFiles(files []interfaces.FileHeader, w http.ResponseWriter) ([]models.VueFile, error)
}

func GetTextContentFromFiles(files []interfaces.FileHeader, w http.ResponseWriter) ([]models.VueFile, error) {
	var fileContents []models.VueFile

	for _, fileHeaderInterface := range files {
		fileHeader := fileHeaderInterface.(*FileHeaderAdapter)
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
		newFile := models.VueFile{Name: fileHeader.FileHeader.Filename, Content: string(fileBytes)}
		fileContents = append(fileContents, newFile)
	}

	return fileContents, nil
}
