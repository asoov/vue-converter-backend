package handlers

import (
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"vue-converter-backend/adapters"
	"vue-converter-backend/cloudwatch"
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
		httpError(w, convertErr, http.StatusBadRequest)
		return
	}

	fileContents, extractErr := s.GetTextContentFromFiles.GetTextContentFromFilesFunc(filesConverted, w)
	if extractErr != nil {
		httpError(w, extractErr, http.StatusBadRequest)
		return
	}

	customerId, customerIdErr := getCustomerId(r)
	if customerIdErr != nil {
		httpError(w, customerIdErr, http.StatusBadRequest)
		return
	}

	customer, customerErr := s.GetCustomer.GetCustomerFunc(customerId)
	if customerErr != nil {
		httpError(w, errors.New("could not retrieve customer"), http.StatusBadRequest)
		return
	}

	if !s.checkTokenBalance(fileContents, customer.AiCredits) {
		httpError(w, errors.New("not enough tokens"), http.StatusBadRequest)
		return
	}

	response, tokensUsed, err := s.GenerateMultipleVueTemplates.GenerateMultipleVueTemplatesFunc(w, r, client, fileContents)
	if err != nil {
		cloudwatch.Log(err.Error())
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	if err := s.deductTokensFromCustomer(customer, tokensUsed); err != nil {
		logAndRespondError(w, "could not deduct tokens", err)
		return
	}

	if err := writeJSONResponse(w, response); err != nil {
		httpError(w, err, http.StatusInternalServerError)
	}
}

func (s *MultipleFiles) parseAndConvertFiles(r *http.Request, w http.ResponseWriter) ([]interfaces.FileHeader, error) {
	files := s.RequestParseFiles.RequestParseFilesFunc(r, w)
	if len(files) == 0 {
		return nil, errors.New("no files uploaded")
	}
	return convertToInterfaceSlice(files), nil
}

func (s *MultipleFiles) deductTokensFromCustomer(customer models.Customer, tokenAmount int) error {
	return s.DeductTokensFromCustomer.DeductTokensForCustomerFunc(customer, tokenAmount)
}

func (s *MultipleFiles) checkTokenBalance(fileContents []models.VueFile, customerAiCreditBalance int) bool {
	calculateNeededTokens := helpers.CalculateNeededTokens{Tokenizer: &interfaces.TokenizerCalToken{}}
	tokensNeeded, err := calculateNeededTokens.CalculateNeededTokensFunc(fileContents)
	if err != nil || customerAiCreditBalance < tokensNeeded {
		return false
	}
	return true
}

func getCustomerId(r *http.Request) (string, error) {
	customerId := r.Header.Get("customer_id")
	if customerId == "" {
		return "", errors.New("no customer id provided in header")
	}
	return customerId, nil
}

func writeJSONResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, writeErr := w.Write(jsonData)
	return writeErr
}

func httpError(w http.ResponseWriter, err error, statusCode int) {
	http.Error(w, err.Error(), statusCode)
}

func logAndRespondError(w http.ResponseWriter, message string, err error) {
	errorMessage := message + ". Reason: " + err.Error()
	cloudwatch.Log(errorMessage)
	http.Error(w, errorMessage, http.StatusInternalServerError)
}

func convertToInterfaceSlice(files []*multipart.FileHeader) []interfaces.FileHeader {
	fileHeaderInterfaceSlice := make([]interfaces.FileHeader, len(files))
	for i, fileHeader := range files {
		fileHeaderInterfaceSlice[i] = adapters.NewFileHeaderAdapter(fileHeader)
	}
	return fileHeaderInterfaceSlice
}
