package services

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"vue-converter-backend/models"

	"github.com/sashabaranov/go-openai"
)

type MockOpenAIClient struct {
	ResponseFunc func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

func (m *MockOpenAIClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	// Call the ResponseFunc if it's defined, otherwise return a default response or error.
	if m.ResponseFunc != nil {
		return m.ResponseFunc(ctx, req)
	}
	return openai.ChatCompletionResponse{}, nil // Return a default response or appropriate error
}

type TestCase struct {
	Name           string
	RequestType    string
	RequestBody    models.GenerateSingleVueTemplateRequest
	Expected       models.GenerateSingleVueTemplateResponse
	ExpectedStatus int
	ClientMock     func(m *MockOpenAIClient)
}

func TestGenerateSingleVueTemplate(t *testing.T) {
	tests := []TestCase{{
		Name:           "Case 1: Filename is empty",
		RequestType:    http.MethodPost,
		RequestBody:    models.GenerateSingleVueTemplateRequest{FileName: "", Content: "This is a test content.", CustomerID: "123"},
		Expected:       models.GenerateSingleVueTemplateResponse{},
		ExpectedStatus: http.StatusBadRequest,
		ClientMock: func(m *MockOpenAIClient) {
			m.ResponseFunc = func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
				return openai.ChatCompletionResponse{
					Usage: openai.Usage{TotalTokens: 10},
					Choices: []openai.ChatCompletionChoice{
						{Message: openai.ChatCompletionMessage{Content: "This is the content of the file."}},
					}}, nil
			}
		},
	},
		{
			Name:           "Case 2: File content is empty",
			RequestType:    http.MethodPost,
			RequestBody:    models.GenerateSingleVueTemplateRequest{FileName: "TestFile.vue", Content: "", CustomerID: "123"},
			Expected:       models.GenerateSingleVueTemplateResponse{},
			ExpectedStatus: http.StatusBadRequest,
			ClientMock:     func(m *MockOpenAIClient) {},
		},
		{
			Name:           "Case 3: OpenAI client returns an error",
			RequestType:    http.MethodPost,
			RequestBody:    models.GenerateSingleVueTemplateRequest{FileName: "TestFile.vue", Content: "<template></template><script></script>", CustomerID: "123"},
			Expected:       models.GenerateSingleVueTemplateResponse{ErrorMessage: "LLM Returned an error"},
			ExpectedStatus: http.StatusInternalServerError,
			ClientMock: func(m *MockOpenAIClient) {
				m.ResponseFunc = func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
					return openai.ChatCompletionResponse{}, &openai.APIError{Message: "LLM Returned an error", Code: 500}
				}
			},
		}, {
			Name:           "Case 4: Returns the File contente and consumed tokens returned by the OpenAI client.",
			RequestType:    http.MethodPost,
			RequestBody:    models.GenerateSingleVueTemplateRequest{FileName: "TestFile.vue", Content: "<template></template><script></script>", CustomerID: "123"},
			Expected:       models.GenerateSingleVueTemplateResponse{FileName: "TestFile.vue", Content: "This is the content of the file.", TokensNeeded: 10},
			ExpectedStatus: http.StatusOK,
			ClientMock: func(m *MockOpenAIClient) {
				m.ResponseFunc = func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
					return openai.ChatCompletionResponse{
						Usage: openai.Usage{TotalTokens: 10},
						Choices: []openai.ChatCompletionChoice{
							{Message: openai.ChatCompletionMessage{Content: "This is the content of the file."}},
						}}, nil
				}
			}},
		{
			Name:        "Case 5: Request is not a POST request",
			RequestType: http.MethodGet,
			RequestBody: models.GenerateSingleVueTemplateRequest{
				FileName:   "TestFile.vue",
				Content:    "<template></template><script></script>",
				CustomerID: "123",
			},
			Expected:       models.GenerateSingleVueTemplateResponse{},
			ExpectedStatus: http.StatusMethodNotAllowed,
			ClientMock:     func(m *MockOpenAIClient) {},
		}}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			mockClient := &MockOpenAIClient{}
			testCase.ClientMock(mockClient)
			marshalledBody, err := json.Marshal(testCase.RequestBody)
			if err != nil {
				t.Error("Error marshalling request body")
			}
			rr := httptest.NewRecorder()
			req, reqCreationError := http.NewRequest(testCase.RequestType, "/generate-web", bytes.NewReader(marshalledBody))
			if reqCreationError != nil {
				t.Error("Error creating request")
			}
			service := &GenerateSingleTemplate{}
			result := service.GenerateSingleVueTemplate(rr, req, mockClient)

			if httpStatus := rr.Code; httpStatus != testCase.ExpectedStatus {
				t.Errorf("HTTP Status is not as expected. Expected: %d, Got: %d", testCase.ExpectedStatus, httpStatus)
			}

			if result.FileName != testCase.Expected.FileName {
				t.Errorf("File name is not as expected. Expected: %s, Got: %s", testCase.Expected.FileName, result.FileName)
			}
			if result.Content != testCase.Expected.Content {
				t.Errorf("Content is not as expected. Expected: %s, Got: %s", testCase.Expected.Content, result.Content)
			}
			if result.TokensNeeded != testCase.Expected.TokensNeeded {
				t.Errorf("Tokens needed is not as expected. Expected: %d, Got: %d", testCase.Expected.TokensNeeded, result.TokensNeeded)
			}
			if result.ErrorMessage != testCase.Expected.ErrorMessage {
				t.Errorf("Error message is not as expected. Expected: %s, Got: %s", testCase.Expected.ErrorMessage, result.ErrorMessage)
			}
		})
	}
}
