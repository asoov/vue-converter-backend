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
	generateSingleResponse GenerateSingleStruct // Injecting the struct directly to avoid having to define every dependency in the main file. It's still possible to provide a different implementation for a different use case.
}

type RequestsWithFileNames struct {
	fileName string
	request  openai.ChatCompletionRequest
}

// Not so nice solution I came along to escape dependency injection hell
// Use a default implementation when no implementation is provided so not EVERY dependency has to be injected in the main file
// If you ahve to inject every dependency in the main file, shit is gonna get extremely bloated and this way you keep flexibility for the cost of some verboseness
// Would be nice to come along another solution for this case

type GenerateSingleStruct struct {
	generateForTest func(request openai.ChatCompletionRequest, client interfaces.OpenAIClient) (openai.ChatCompletionResponse, error)
}

func (s *GenerateSingleStruct) GenerateSingleTemplateResponse(request openai.ChatCompletionRequest, client interfaces.OpenAIClient) (openai.ChatCompletionResponse, error) {
	if s.generateForTest != nil {
		return s.generateForTest(request, client)
	}
	generateSingleFunc := GenerateSingleTemplateResponse{}
	return generateSingleFunc.GenerateSingleTemplateResponseFunc(request, client)
}

func (s *GenerateMultipleVueTemplates) GenerateMultipleVueTemplatesFunc(w http.ResponseWriter, r *http.Request, client interfaces.OpenAIClient, files []models.VueFile) (models.GenerateMultipleVueTemplateResponse, error) {
	requests := createRequestsWithFileNames(files)

	if len(requests) == 0 {
		return models.GenerateMultipleVueTemplateResponse{}, errors.New("no files uploaded")
	}

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
			println("Processing request for file: ", req.fileName)

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
