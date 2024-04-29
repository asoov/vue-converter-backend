package services

import (
	"errors"
	"net/http"
	"vue-converter-backend/helpers"
	"vue-converter-backend/interfaces"
	"vue-converter-backend/models"

	"github.com/sashabaranov/go-openai"
)

type GenerateMultipleVueTemplates struct {
	generateSingleResponse GenerateSingleTemplateResponseInterface
}

type Files struct {
	fileName    string
	fileContent string
}

type RequestsWithFileNames struct {
	fileName string
	request  openai.ChatCompletionRequest
}

func (s *GenerateMultipleVueTemplates) GenerateMultipleVueTemplatesFunc(w http.ResponseWriter, r *http.Request, client interfaces.OpenAIClient, files []Files) (*models.GenerateMultipleVueTemplateResponse, error) {
	requests := []RequestsWithFileNames{}
	for _, fileContent := range files {
		chatRequest := GetChatRequest(fileContent.fileContent)
		requests = append(requests, RequestsWithFileNames{fileName: fileContent.fileName, request: chatRequest})
	}

	if len(requests) == 0 {
		return nil, errors.New("no files uploaded")
	}

	var results models.GenerateMultipleVueTemplateResponse

	for _, request := range requests {
		chatCompletionResult, err := s.generateSingleResponse.GenerateSingleTemplateResponse(request.request, client)
		mappedResult := helpers.MapOpenAiResponse(request.fileName, chatCompletionResult, err)
		results = append(results, mappedResult)
	}

	return &results, nil
}
