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

func (f *FileHeaderAdapter) Filename() string {
	return f.FileHeader.Filename
}

type GetTextContentFromFileInterface interface {
	GetTextContentFromFilesFunc(files []interfaces.FileHeader, w http.ResponseWriter) ([]models.VueFile, error)
}

type GetTextContentFromFiles struct{}

func (*GetTextContentFromFiles) GetTextContentFromFilesFunc(files []interfaces.FileHeader, w http.ResponseWriter) ([]models.VueFile, error) {
	var fileContents []models.VueFile

	for _, fileHeaderInterface := range files {
		fileHeader := fileHeaderInterface

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
		fileName := fileHeader.Filename()
		newFile := models.VueFile{Name: fileName, Content: string(fileBytes)}
		fileContents = append(fileContents, newFile)
	}

	return fileContents, nil
}
