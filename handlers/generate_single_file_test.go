package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"vue-converter-backend/interfaces"
	"vue-converter-backend/models"
	"vue-converter-backend/services"

	"github.com/sashabaranov/go-openai"
)

type MockOpenAIClient struct {
	createChatCompletionReplacement func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

func (s *MockOpenAIClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	if s.createChatCompletionReplacement != nil {
		return s.createChatCompletionReplacement(ctx, req)
	}
	return openai.NewClient(os.Getenv("OAI_KEY")).CreateChatCompletion(ctx, req)
}

type MockGenerateSingleResponse struct {
	generateSingleResponseReplacement func(w http.ResponseWriter, r *http.Request, client interfaces.OpenAIClient) models.GenerateSingleVueTemplateResponse
}

func (s *MockGenerateSingleResponse) GenerateSingleVueTemplateFunc(w http.ResponseWriter, r *http.Request, client interfaces.OpenAIClient) models.GenerateSingleVueTemplateResponse {
	if s.generateSingleResponseReplacement != nil {
		return s.generateSingleResponseReplacement(w, r, client)
	}

	generateSingleResponse := services.GenerateSingleTemplate{}

	return generateSingleResponse.GenerateSingleVueTemplateFunc(w, r, client)
}

type testCase struct {
	name                  string
	mocks                 func(o *MockOpenAIClient, g *MockGenerateSingleResponse)
	given                 functionInput
	expectedSuccessResult models.GenerateSingleVueTemplateResponse
	expectedErrorMessage  string
}

type functionInput struct {
	request *http.Request
	client  interfaces.OpenAIClient
}

func TestGenerateSingleFile(t *testing.T) {
	testCases := []testCase{
		{
			name: "Test Case 1: It should return an error if OAI client is not initialized",
			given: functionInput{
				request: httptest.NewRequest(http.MethodPost, "/generate-single-file", nil),
				client:  nil,
			},
			mocks:                func(o *MockOpenAIClient, g *MockGenerateSingleResponse) {},
			expectedErrorMessage: "OpenAI client is not initialized",
		},
		{
			name: "Test Case 2: It should return an error if JSON marshalling fails",
			given: functionInput{
				request: httptest.NewRequest(http.MethodPost, "/generate-single-file", nil),
				client:  openai.NewClient(os.Getenv("OAI_KEY")),
			},
			mocks: func(o *MockOpenAIClient, g *MockGenerateSingleResponse) {
				o.createChatCompletionReplacement = func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
					return openai.ChatCompletionResponse{}, nil
				}
				g.generateSingleResponseReplacement = func(w http.ResponseWriter, r *http.Request, client interfaces.OpenAIClient) models.GenerateSingleVueTemplateResponse {
					// Create invalid UTF8 to make JSON marshalling fail
					return models.GenerateSingleVueTemplateResponse{Content: }
				}
			},
			expectedErrorMessage: "json: error calling MarshalJSON for type json.RawMessage: invalid character '\\xff' looking for beginning of value\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			openaiClient := &MockOpenAIClient{}
			generateSingleMock := &MockGenerateSingleResponse{}
			tc.mocks(openaiClient, generateSingleMock)
			generateSingleFile := GenerateSingleFile{GenerateSingleFile: generateSingleMock}

			responseWriter := httptest.NewRecorder()

			generateSingleFile.GenerateSingleFileFunc(responseWriter, tc.given.request, tc.given.client)

			// Assertions for the failing cases where the error code will be written to the response body
			expectedErrorMessage := strings.TrimSpace(tc.expectedErrorMessage)
			actualErrorMessage := strings.TrimSpace(responseWriter.Body.String())
			if responseWriter.Code != http.StatusOK && expectedErrorMessage != actualErrorMessage {
				t.Errorf("Expected error message: %v, but got %v", expectedErrorMessage, actualErrorMessage)
			}

			// Assertions for the successful cases
			if responseWriter.Code == http.StatusOK {
				var response models.GenerateSingleVueTemplateResponse
				err := json.Unmarshal(responseWriter.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Error unmarshalling response: %v", err)
				}
				if !reflect.DeepEqual(response, tc.expectedSuccessResult) {
					t.Errorf("Expected response: %v, but got %v", tc.expectedSuccessResult, response)
				}
			}
		})
	}
}
