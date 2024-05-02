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

type RequestsWithFileNames struct {
	fileName string
	request  openai.ChatCompletionRequest
}

func (s *GenerateMultipleVueTemplates) GenerateMultipleVueTemplatesFunc(w http.ResponseWriter, r *http.Request, client interfaces.OpenAIClient, files []models.VueFile) (models.GenerateMultipleVueTemplateResponse, error) {
	requests := []RequestsWithFileNames{}
	for _, file := range files {
		chatRequest := GetChatRequest(file.Content)
		requests = append(requests, RequestsWithFileNames{fileName: file.Name, request: chatRequest})
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

	return results, nil
}
