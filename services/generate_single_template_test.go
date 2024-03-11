package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"vue-converter-backend/models"

	"github.com/sashabaranov/go-openai"
)

type MockOpenAIClient struct{}

func (m *MockOpenAIClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	// Return a mock response or error as needed for tests
	return openai.ChatCompletionResponse{
		ID: "",
		Usage: openai.Usage{
			TotalTokens: 10,
		},
		Choices: []openai.ChatCompletionChoice{{
			Message: openai.ChatCompletionMessage{
				Content: "This is a generated Completion message.",
			},
		}},
	}, nil
}

func TestGenerateSingleVueTemplate(t *testing.T) {
	mockClient := &MockOpenAIClient{}
	requestBody := models.GenerateSingleVueTemplateRequest{FileName: "test.vue", Content: "This is a test content."}
	marshalledBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Error("Error marshalling request body")
	}
	rr := httptest.NewRecorder()
	req, reqCreationError := http.NewRequest(http.MethodPost, "/generate-web", bytes.NewReader(marshalledBody))
	if reqCreationError != nil {
		t.Error("Error creating request")
	}
	result := GenerateSingleVueTemplate(rr, req, mockClient)
	println("Hallo")
	fmt.Printf(result.Content)

	if result.FileName != "test.vue" {
		t.Error("File name is not as expected")
	}
	if result.TokensNeeded != 10 {
		t.Errorf("Tokens needed is not as expected.")
	}

}
