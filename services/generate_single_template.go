package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"vue-converter-backend/interfaces"
	"vue-converter-backend/models"
)

type RegexpCompile struct{}

func (m *RegexpCompile) Compile(str string) (*regexp.Regexp, error) {
	return regexp.MustCompile(str), nil
}

type GenerateSingleVueTemplateInterface interface {
	GenerateSingleVueTemplateFunc(w http.ResponseWriter, r *http.Request, client interfaces.OpenAIClient) models.GenerateSingleVueTemplateResponse
}

type GenerateSingleTemplate struct{}

func (s *GenerateSingleTemplate) GenerateSingleVueTemplateFunc(w http.ResponseWriter, r *http.Request, client interfaces.OpenAIClient) models.GenerateSingleVueTemplateResponse {
	// Only process POST requests
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return models.GenerateSingleVueTemplateResponse{}
	}
	// Step 2: Read the request body
	body, readErr := readRequestBody(r.Body, w)

	if readErr != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return models.GenerateSingleVueTemplateResponse{}
	}

	var requestBody models.GenerateSingleVueTemplateRequest

	parsingError := json.Unmarshal(body, &requestBody)

	if parsingError != nil {
		fmt.Println(parsingError)
		http.Error(w, "Error parsing request", http.StatusBadRequest)
		return models.GenerateSingleVueTemplateResponse{}
	}
	if requestBody.FileName == "" || requestBody.Content == "" {
		http.Error(w, "Error parsing request", http.StatusBadRequest)
		return models.GenerateSingleVueTemplateResponse{}
	}

	chatCompletionRequest := GetChatRequest(requestBody.Content)
	handler := GenerateSingleTemplateResponse{}
	result, generationError := handler.GenerateSingleTemplateResponseFunc(chatCompletionRequest, client)
	var errorMessage string = ""
	if generationError != nil {
		errorMessage = generationError.Error()
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return models.GenerateSingleVueTemplateResponse{ErrorMessage: errorMessage}
	}

	fileContentResult := result.Choices[0].Message.Content
	tokensConsumed := result.Usage.TotalTokens

	return models.GenerateSingleVueTemplateResponse{
		FileName:     requestBody.FileName,
		Content:      fileContentResult,
		TokensNeeded: tokensConsumed,
		ErrorMessage: errorMessage,
	}
}

func readRequestBody(requestBody io.ReadCloser, w http.ResponseWriter) ([]byte, error) {
	body, err := io.ReadAll(requestBody)

	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}
	defer requestBody.Close()
	return body, err
}
