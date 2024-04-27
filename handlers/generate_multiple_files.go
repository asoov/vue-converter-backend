package handlers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"vue-converter-backend/adapters"
	"vue-converter-backend/helpers"
	"vue-converter-backend/interfaces"

	"github.com/sashabaranov/go-openai"
)

type MultipleFiles struct {
	generate interfaces.GenerateMultipleVueTemplates
}

func (s *MultipleFiles) GenerateMultipleFiles(w http.ResponseWriter, r *http.Request, client *openai.Client) {

	files := helpers.RequestParseFiles(r, w)

	filesConverted := convertToInterfaceSlice(files)

	fileContents, err := helpers.GetTextContentFromFiles(filesConverted, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := s.generate.GenerateMultipleVueTemplates(w, r, client, fileContents)
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(response)

	if err != nil {
		// If an error occurs during JSON marshaling, send an error to the client
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
