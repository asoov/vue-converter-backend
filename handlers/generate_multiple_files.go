package handlers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"vue-converter-backend/adapters"
	"vue-converter-backend/helpers"
	"vue-converter-backend/interfaces"
	"vue-converter-backend/services"

	"github.com/sashabaranov/go-openai"
)

type MultipleFiles struct {
	GenerateMultipleVueTemplates services.GenerateMultipleInterface
	GetTextContentFromFiles      helpers.GetTextContentFromFileInterface
	RequestParseFiles            helpers.RequestParseFilesInterface
}

func (s *MultipleFiles) GenerateMultipleFilesFunc(w http.ResponseWriter, r *http.Request, client *openai.Client) {

	files := s.RequestParseFiles.RequestParseFilesFunc(r, w)

	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	filesConverted := convertToInterfaceSlice(files)

	fileContents, extractErr := s.GetTextContentFromFiles.GetTextContentFromFiles(filesConverted, w)

	if extractErr != nil {
		http.Error(w, extractErr.Error(), http.StatusBadRequest)
		return
	}

	response := s.GenerateMultipleVueTemplates.GenerateMultipleVueTemplatesFunc(w, r, client, fileContents)

	w.Header().Set("Content-Type", "application/json")

	jsonData, marshalErr := json.Marshal(response)

	if marshalErr != nil {
		// If an error occurs during JSON marshaling, send an error to the client
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}
	_, writeErr := w.Write(jsonData)
	if writeErr != nil {
		// If there is an error while writing, log it, or handle it as necessary
		http.Error(w, "Failed to write JSON to response", http.StatusInternalServerError)
		return
	}
}

func convertToInterfaceSlice(files []*multipart.FileHeader) []interfaces.FileHeader {
	var fileHeaderInterfaceSlice []interfaces.FileHeader

	for _, fileHeader := range files {
		fileHeaderInterfaceSlice = append(fileHeaderInterfaceSlice, adapters.NewFileHeaderAdapter(fileHeader))
	}

	return fileHeaderInterfaceSlice
}
