package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"vue-converter-backend/helpers"
	"vue-converter-backend/models"

	"github.com/sashabaranov/go-openai"
)

func GenerateSingleVueTemplate(w http.ResponseWriter, r *http.Request) models.GenerateSingleVueTemplateResponse {
	// Only process POST requests
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
	}
	// Step 2: Read the request body
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
	}
	defer r.Body.Close()

	var requestBody models.GenerateSingleVueTemplateRequest

	parsingError := json.Unmarshal(body, &requestBody)

	if parsingError != nil {
		fmt.Println(parsingError)
		http.Error(w, "Error parsing request", http.StatusInternalServerError)
		return models.GenerateSingleVueTemplateResponse{}
	}

	neededTokens := helpers.CalculateNeededTokens(requestBody.Content, w)

	result, generationError := generateSingleTemplateResponse(requestBody.Content, neededTokens, w)

	fileContentResult := result.Choices[0].Message.Content
	tokensConsumed := result.Usage.TotalTokens

	var errorMessage string
	if generationError != nil {
		errorMessage = generationError.Error()
	}

	return models.GenerateSingleVueTemplateResponse{
		FileName:     requestBody.FileName,
		Content:      fileContentResult,
		TokensNeeded: tokensConsumed,
		ErrorMessage: &errorMessage,
	}
}

func generateSingleTemplateResponse(fileContent string, neededTokens int, w http.ResponseWriter) (openai.ChatCompletionResponse, error) {
	client := openai.NewClient(os.Getenv("OAI_KEY"))

	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo-16k",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "user",
				Content: GetVueChatCompletion(fileContent),
			},
		},
	})

	return resp, err
}
