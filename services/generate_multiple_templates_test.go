package services

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"
	"vue-converter-backend/interfaces"
	"vue-converter-backend/models"

	"github.com/sashabaranov/go-openai"
)

type MockOAI struct{}

type GenerateSingleResponseStruct struct{}

func (s *MockOAI) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	return openai.NewClient(os.Getenv("OAI_KEY")).CreateChatCompletion(ctx, req)
}

func (s *GenerateSingleResponseStruct) GenerateSingleTemplateResponse(request openai.ChatCompletionRequest, client interfaces.OpenAIClient) (openai.ChatCompletionResponse, error) {
	println("DRIN")
	return client.CreateChatCompletion(context.Background(), request)
}

func TestGenerateMultipleTemplates(t *testing.T) {
	type TestCase struct {
		name  string
		files []models.VueFile
		mock  func(*GenerateSingleResponseStruct)
	}

	testCases := []TestCase{{
		name:  "It should not call the OpenAI API if the file content is empty",
		files: []models.VueFile{{Name: "test.vue", Content: ""}},
		mock:  func(*GenerateSingleResponseStruct) {},
	}}

	for _, tc := range testCases {
		tc.name = "Test case 1"

		handler := GenerateMultipleVueTemplates{generateSingleResponse: &GenerateSingleResponseStruct{}}

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/generate-web/multiple", nil)
		mock := &MockOAI{}
		handler.GenerateMultipleVueTemplatesFunc(rr, req, mock, tc.files)
	}
}
