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

type MockOAI struct {
}

type MockGenerateSingleResponse struct {
	mockReturn func(request openai.ChatCompletionRequest, client interfaces.OpenAIClient) (openai.ChatCompletionResponse, error)
}

func (s *MockOAI) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	return openai.NewClient(os.Getenv("OAI_KEY")).CreateChatCompletion(ctx, req)
}

func (s *MockGenerateSingleResponse) GenerateSingleTemplateResponse(request openai.ChatCompletionRequest, client interfaces.OpenAIClient) (openai.ChatCompletionResponse, error) {
	return s.mockReturn(request, client)

}

func TestGenerateMultipleTemplates(t *testing.T) {
	type TestCase struct {
		name                       string
		files                      []models.VueFile
		generateSingleResponseMock func(*MockGenerateSingleResponse)
		expected                   models.GenerateMultipleVueTemplateResponse
		expectedErrorMessage       string
	}

	testCases := []TestCase{{
		name:  "It should return an error if files are empty",
		files: []models.VueFile{},
		generateSingleResponseMock: func(m *MockGenerateSingleResponse) {
			m.mockReturn = func(request openai.ChatCompletionRequest, client interfaces.OpenAIClient) (openai.ChatCompletionResponse, error) {
				return openai.ChatCompletionResponse{}, nil
			}
		},
		expected:             nil,
		expectedErrorMessage: "no files uploaded",
	}, {
		name:  "It should return an array of chat completion responses",
		files: []models.VueFile{{Name: "Name", Content: "This is the content"}},
		generateSingleResponseMock: func(m *MockGenerateSingleResponse) {
			m.mockReturn = func(request openai.ChatCompletionRequest, client interfaces.OpenAIClient) (openai.ChatCompletionResponse, error) {
				return openai.ChatCompletionResponse{
						Choices: []openai.ChatCompletionChoice{
							openai.ChatCompletionChoice{
								Message: openai.ChatCompletionMessage{Content: "Hallo"},
							},
						},
						Usage: openai.Usage{TotalTokens: 2},
					},
					nil
			}
		},
		expected:             models.GenerateMultipleVueTemplateResponse{models.GenerateSingleVueTemplateResponse{FileName: "Name", Content: "Hallo", TokensNeeded: 2}},
		expectedErrorMessage: "",
	}}

	for _, tc := range testCases {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/generate-web/multiple", nil)
		openAiMock := &MockOAI{}
		generateSingleResponseMock := &MockGenerateSingleResponse{}
		tc.generateSingleResponseMock(generateSingleResponseMock)

		handler := GenerateMultipleVueTemplates{generateSingleResponse: GenerateSingleStruct{generateForTest: generateSingleResponseMock.GenerateSingleTemplateResponse}}
		result, error := handler.GenerateMultipleVueTemplatesFunc(rr, req, openAiMock, tc.files)

		if error == nil {
			if result[0] != tc.expected[0] {
				t.Errorf("Result is not what was expected, expected: %v got: %v", tc.expected, result)
			}
		} else {
			if error.Error() != tc.expectedErrorMessage {
				t.Error("Error not what was expected")
			}
		}
	}
}
