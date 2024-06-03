package services

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"strings"
	"vue-converter-backend/models"
)

type CreateZipFile struct{}

type CreateZipFileInterface interface {
	CreateZipFileFunc(files models.GenerateMultipleVueTemplateResponse) ([]byte, error)
}

func (s *CreateZipFile) CreateZipFileFunc(files models.GenerateMultipleVueTemplateResponse) ([]byte, error) {

	var buffer bytes.Buffer
	zipWriter := zip.NewWriter(&buffer)

	defer zipWriter.Close()

	for _, file := range files {
		zipFile, err := zipWriter.Create(file.FileName)
		if err != nil {
			log.Fatal(err)
		}

		fileContentAsReader := strings.NewReader(file.Content)

		_, err = io.Copy(zipFile, fileContentAsReader)
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil

}
