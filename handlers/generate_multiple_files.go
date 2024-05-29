package handlers

import (
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"vue-converter-backend/adapters"
	"vue-converter-backend/dynamo"
	"vue-converter-backend/helpers"
	"vue-converter-backend/interfaces"
	"vue-converter-backend/models"

	"github.com/sashabaranov/go-openai"
)

type MultipleFiles struct {
	GenerateMultipleVueTemplates interfaces.GenerateMultipleVueTemplates
	GetTextContentFromFiles      helpers.GetTextContentFromFileInterface
	RequestParseFiles            helpers.RequestParseFilesInterface
	GetCustomer                  dynamo.GetCustomerInterface
	DeductTokensFromCustomer     dynamo.DeductTokensForCustomerInterface
}

func (s *MultipleFiles) GenerateMultipleFilesFunc(w http.ResponseWriter, r *http.Request, client *openai.Client) {

	filesConverted, convertErr := s.parseAndConvertFiles(r, w)

	if convertErr != nil {
		http.Error(w, convertErr.Error(), http.StatusBadRequest)
		return
	}

	fileContents, extractErr := s.GetTextContentFromFiles.GetTextContentFromFilesFunc(filesConverted, w)

	if extractErr != nil {
		http.Error(w, extractErr.Error(), http.StatusBadRequest)
		return
	}

	customerId := r.Header.Get("customer_id")

	if customerId == "" {
		http.Error(w, "No customer id provided in header", http.StatusBadRequest)
		return
	}

	customer, customerErr := s.GetCustomer.GetCustomerFunc(customerId)
	if customerErr != nil {
		http.Error(w, "Could not retrieve customer", http.StatusBadRequest)
		return
	}

	ok := s.checkTokenBalance(fileContents, customer.AiCredits)

	if !ok {
		http.Error(w, "Not enough tokens", http.StatusBadRequest)
		return
	}

	response, tokensUsed, err := s.GenerateMultipleVueTemplates.GenerateMultipleVueTemplatesFunc(w, r, client, fileContents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.deductTokensFromCustomer(customer, tokensUsed)

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

func (s *MultipleFiles) parseAndConvertFiles(r *http.Request, w http.ResponseWriter) ([]interfaces.FileHeader, error) {
	files := s.RequestParseFiles.RequestParseFilesFunc(r, w)

	if len(files) == 0 {
		return nil, errors.New("no files uploaded")
	}

	filesConverted := convertToInterfaceSlice(files)
	return filesConverted, nil
}

func (s *MultipleFiles) deductTokensFromCustomer(customer models.Customer, tokenAmount int) error {
	return s.DeductTokensFromCustomer.DeductTokensForCustomerFunc(customer, tokenAmount)
}

func (s *MultipleFiles) checkTokenBalance(fileContents []models.VueFile, customerAiCreditBalance int) bool {

	calculateNeededTokens := helpers.CalculateNeededTokens{Tokenizer: &interfaces.TokenizerCalToken{}}
	tokensNeeded, err := calculateNeededTokens.CalculateNeededTokensFunc(fileContents)
	if err != nil {
		return false
	}

	if customerAiCreditBalance < tokensNeeded {
		return false
	} else {
		return true
	}
}

func (s *MultipleFiles) GetCustomerFunc(id string) (models.Customer, error) {
	return s.GetCustomer.GetCustomerFunc(id)
}
