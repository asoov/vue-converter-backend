package services

import (
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
	requests := createRequestsWithFileNames(files)

	const maxConcurrentRequests = 10

	resultsChannel := make(chan models.GenerateSingleVueTemplateResponse, len(requests))

	processRequestsConcurrently(requests, resultsChannel, s.generateSingleResponse.GenerateSingleTemplateResponse, client)

	results := collectResults(len(requests), resultsChannel)
	return results, nil
}

func createRequestsWithFileNames(files []models.VueFile) []RequestsWithFileNames {
	requests := []RequestsWithFileNames{}
	for _, file := range files {
		chatRequest := GetChatRequest(file.Content)
		requests = append(requests, RequestsWithFileNames{fileName: file.Name, request: chatRequest})
	}
	return requests
}

func processRequestsConcurrently(requests []RequestsWithFileNames, resultsChannel chan models.GenerateSingleVueTemplateResponse, generateSingleTemplateResponse func(request openai.ChatCompletionRequest, client interfaces.OpenAIClient) (openai.ChatCompletionResponse, error), client interfaces.OpenAIClient) {
	// Limit the number of concurrent requests to 10
	semaphore := make(chan struct{}, 10)

	for _, request := range requests {
		go func(req RequestsWithFileNames) {
			semaphore <- struct{}{}

			chatCompletionResult, err := generateSingleTemplateResponse(req.request, client)
			mappedResult := helpers.MapOpenAiResponse(req.fileName, chatCompletionResult, err)
			resultsChannel <- mappedResult

			<-semaphore
		}(request)
	}
}

func collectResults(length int, resultsChannel chan models.GenerateSingleVueTemplateResponse) models.GenerateMultipleVueTemplateResponse {
	var results models.GenerateMultipleVueTemplateResponse
	for i := 0; i < length; i++ {
		result := <-resultsChannel
		results = append(results, result)
	}
	return results
}
