package handlers

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"vue-converter-backend/dynamo"
	"vue-converter-backend/helpers"
	"vue-converter-backend/interfaces"
	"vue-converter-backend/models"

	"github.com/sashabaranov/go-openai"
)

type MockGenerateMultiple struct {
	execute func(w http.ResponseWriter, r *http.Request, client interfaces.OpenAIClient, files []models.VueFile) models.GenerateMultipleVueTemplateResponse
}

type MockParseFiles struct {
	execParseFiles func(r *http.Request, w http.ResponseWriter) []*multipart.FileHeader
}

type MockGetTextContent struct {
	execGetTextContent func(files []interfaces.FileHeader, w http.ResponseWriter) ([]models.VueFile, error)
}

type MockGetCustomer struct {
	getCustomerReplacement func(id string) (models.Customer, error)
}

type MockDeductTokensForCustomer struct{}

func (s *MockGenerateMultiple) GenerateMultipleVueTemplatesFunc(w http.ResponseWriter, r *http.Request, client interfaces.OpenAIClient, files []models.VueFile) (models.GenerateMultipleVueTemplateResponse, int, error) {
	if s.execute != nil {
		result := s.execute(w, r, client, files)
		return result, 0, nil
	}
	return models.GenerateMultipleVueTemplateResponse{}, 0, nil
}

func (s *MockParseFiles) RequestParseFilesFunc(r *http.Request, w http.ResponseWriter) []*multipart.FileHeader {
	if s.execParseFiles != nil {
		return s.execParseFiles(r, w)
	}
	handler := helpers.RequestParseFiles{}
	return handler.RequestParseFilesFunc(r, w)
}
func (s *MockGetTextContent) GetTextContentFromFilesFunc(files []interfaces.FileHeader, w http.ResponseWriter) ([]models.VueFile, error) {
	if s.execGetTextContent != nil {
		return s.execGetTextContent(files, w)
	}
	getTextContentFromFiles := helpers.GetTextContentFromFiles{}
	return getTextContentFromFiles.GetTextContentFromFilesFunc(files, w)
}

func (s *MockGetCustomer) GetCustomerFunc(id string) (models.Customer, error) {
	if s.getCustomerReplacement != nil {
		return s.getCustomerReplacement(id)
	}
	getCustomerFunc := dynamo.GetCustomer{}

	return getCustomerFunc.GetCustomerFunc(id)
}

func (s *MockDeductTokensForCustomer) DeductTokensForCustomerFunc(customer models.Customer, tokenAmount int) error {
	return nil

}

type TestCase struct {
	name       string
	mock       func(s *MockGenerateMultiple, d *MockGetTextContent, f *MockParseFiles, g *MockGetCustomer)
	givenFiles []File
	assertion  func(r *httptest.ResponseRecorder, t *testing.T)
}

type File struct {
	filename string
	content  string
}

func TestGenerateMultipleFiles(t *testing.T) {
	testCases := []TestCase{
		{
			name:       "Case 1: If request contains no files, return a 400 status code",
			mock:       func(s *MockGenerateMultiple, d *MockGetTextContent, f *MockParseFiles, g *MockGetCustomer) {},
			givenFiles: []File{},
			assertion: func(r *httptest.ResponseRecorder, t *testing.T) {
				if r.Code != http.StatusBadRequest {
					t.Error("Not a 400 error code")
				}
				if strings.TrimSpace(r.Body.String()) != "no files uploaded" {
					t.Errorf("Wrong message returned %s", r.Body.String())
				}
			},
		},
		{
			name: "Case 2: When request goes through, header should be set, correct status code returned, and contetn should be correct",
			mock: func(s *MockGenerateMultiple, d *MockGetTextContent, f *MockParseFiles, g *MockGetCustomer) {
				g.getCustomerReplacement = func(id string) (models.Customer, error) {
					return models.Customer{FirstName: "Anton", AiCredits: 120000}, nil
				}
				s.execute = func(w http.ResponseWriter, r *http.Request, client interfaces.OpenAIClient, files []models.VueFile) models.GenerateMultipleVueTemplateResponse {
					return models.GenerateMultipleVueTemplateResponse{models.GenerateSingleVueTemplateResponse{
						FileName:     "NameName",
						Content:      "This is the content",
						TokensNeeded: 666,
					}}
				}
			},
			givenFiles: []File{{filename: "fileName1", content: "Diese"}, {filename: "fileName2", content: "Das ist ein Test"}},
			assertion: func(r *httptest.ResponseRecorder, t *testing.T) {
				if r.Code != http.StatusOK {
					t.Errorf("Doesnt have correct status code, expected: %v, is: %v", http.StatusOK, r.Code)
				}
				if r.Header().Get("Content-Type") != "application/json" {
					t.Error("Doenst have application/json header")
				}
			},
		},
		{
			name: "Case 3: Create error when json marshaling fails",
			mock: func(s *MockGenerateMultiple, d *MockGetTextContent, f *MockParseFiles, g *MockGetCustomer) {
				s.execute = func(w http.ResponseWriter, r *http.Request, client interfaces.OpenAIClient, files []models.VueFile) models.GenerateMultipleVueTemplateResponse {
					// TODO: Create a response that will break json.Marshal()
					return models.GenerateMultipleVueTemplateResponse{models.GenerateSingleVueTemplateResponse{Content: ""}}
				}
			},
		},
		{
			name: "Case 4: Write json response data.",
			mock: func(s *MockGenerateMultiple, d *MockGetTextContent, f *MockParseFiles, g *MockGetCustomer) {
				s.execute = func(w http.ResponseWriter, r *http.Request, client interfaces.OpenAIClient, files []models.VueFile) models.GenerateMultipleVueTemplateResponse {
					// TODO: Create a response that will break json.Marshal()
					return models.GenerateMultipleVueTemplateResponse{models.GenerateSingleVueTemplateResponse{Content: ""}}
				}
			},
		},
		{
			name: "Case 5: Error when extraction fails",
			mock: func(s *MockGenerateMultiple, d *MockGetTextContent, f *MockParseFiles, g *MockGetCustomer) {
				d.execGetTextContent = func(files []interfaces.FileHeader, w http.ResponseWriter) ([]models.VueFile, error) {
					return nil, errors.New("Error occured")
				}
			},
			givenFiles: []File{{filename: "fileName1", content: "Diese"}, {filename: "fileName2", content: "Das ist ein Test"}},
			assertion: func(r *httptest.ResponseRecorder, t *testing.T) {
				if r.Code != http.StatusBadRequest {
					t.Errorf("Wrong code! expected: %v got: %v", http.StatusBadRequest, r.Code)
				}

				expectedErrString := "Error occured"
				returnedStringTrimmed := strings.TrimSpace(r.Body.String())
				if returnedStringTrimmed != expectedErrString {
					t.Errorf("Wrong error! expected: %s got: %s", expectedErrString, strings.TrimSpace(r.Body.String()))
				}
			},
		},
		{
			name: "Case 6: Error when getting customer fails",
			mock: func(s *MockGenerateMultiple, d *MockGetTextContent, f *MockParseFiles, g *MockGetCustomer) {
				g.getCustomerReplacement = func(id string) (models.Customer, error) {
					return models.Customer{}, errors.New("Error occured")
				}
			},
			givenFiles: []File{{filename: "fileName1", content: "Diese"}, {filename: "fileName2", content: "Das ist ein Test"}},
			assertion: func(r *httptest.ResponseRecorder, t *testing.T) {
				if r.Code != http.StatusBadRequest {
					t.Errorf("Wrong code! expected: %v got: %v", http.StatusBadRequest, r.Code)
				}

				expectedErrString := strings.TrimSpace("Could not retrieve customer")
				returnedStringTrimmed := strings.TrimSpace(r.Body.String())
				if returnedStringTrimmed != expectedErrString {
					t.Errorf("Wrong error! expected: %s got: %s", expectedErrString, returnedStringTrimmed)
				}
			},
		},
		{
			name: "Case 7: User doesnt have enough tokens",
			mock: func(s *MockGenerateMultiple, d *MockGetTextContent, f *MockParseFiles, g *MockGetCustomer) {
				g.getCustomerReplacement = func(id string) (models.Customer, error) {
					return models.Customer{AiCredits: 0}, nil
				}
			},
			givenFiles: []File{{filename: "fileName1", content: "Diese"}, {filename: "fileName2", content: "Das ist ein Test"}},
			assertion: func(r *httptest.ResponseRecorder, t *testing.T) {
				if r.Code != http.StatusBadRequest {
					t.Errorf("Wrong code! expected: %v got: %v", http.StatusBadRequest, r.Code)
				}

				expectedErrString := strings.TrimSpace("Not enough tokens")
				returnedStringTrimmed := strings.TrimSpace(r.Body.String())

				if returnedStringTrimmed != expectedErrString {
					t.Errorf("Wrong error! expected: %s got: %s", expectedErrString, returnedStringTrimmed)
				}
			},
		},
		{
			name: "Case 8: User has enough tokens",
			mock: func(s *MockGenerateMultiple, d *MockGetTextContent, f *MockParseFiles, g *MockGetCustomer) {
				g.getCustomerReplacement = func(id string) (models.Customer, error) {
					return models.Customer{AiCredits: 99999}, nil
				}
				s.execute = func(w http.ResponseWriter, r *http.Request, client interfaces.OpenAIClient, files []models.VueFile) models.GenerateMultipleVueTemplateResponse {
					return models.GenerateMultipleVueTemplateResponse{models.GenerateSingleVueTemplateResponse{
						FileName:     "NameName",
						Content:      "This is the content",
						TokensNeeded: 666,
					}}
				}
			},
			givenFiles: []File{{filename: "fileName1", content: "Diese"}, {filename: "fileName2", content: "Das ist ein Test"}},
			assertion: func(r *httptest.ResponseRecorder, t *testing.T) {
				if r.Code != http.StatusOK {
					t.Errorf("Wrong code! expected: %v got: %v", http.StatusOK, r.Code)
				}
			},
		},
	}

	for _, tc := range testCases {
		mockGenerateMultiple := MockGenerateMultiple{}
		mockParseFiles := MockParseFiles{}
		mockGetTextContent := MockGetTextContent{}
		mockGetCustomer := MockGetCustomer{}
		mockDeductTokensForCustomer := MockDeductTokensForCustomer{}

		tc.mock(&mockGenerateMultiple, &mockGetTextContent, &mockParseFiles, &mockGetCustomer)

		multipleFiles := MultipleFiles{
			GenerateMultipleVueTemplates: &mockGenerateMultiple,
			GetTextContentFromFiles:      &mockGetTextContent,
			RequestParseFiles:            &mockParseFiles,
			GetCustomer:                  &mockGetCustomer,
			DeductTokensFromCustomer:     &mockDeductTokensForCustomer,
		}

		rr := httptest.NewRecorder()
		req, err := createMultifileUploadRequest(tc.givenFiles)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}

		client := openai.NewClient("your-api-key")
		multipleFiles.GenerateMultipleFilesFunc(rr, req, client)
		if tc.assertion != nil {
			tc.assertion(rr, t)
		}
	}
}

func createMultifileUploadRequest(files []File) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, file := range files {
		part, createErr := writer.CreateFormFile("files", file.filename)
		if createErr != nil {
			return nil, createErr
		}

		_, writeErr := part.Write([]byte(file.content))
		if writeErr != nil {
			return nil, writeErr
		}
	}

	closeErr := writer.Close()
	if closeErr != nil {
		return nil, closeErr
	}

	req := httptest.NewRequest(http.MethodPost, "/", body)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("customer_id", "1234")

	return req, nil
}
