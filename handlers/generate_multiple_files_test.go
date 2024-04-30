package handlers

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

func (s *MockGenerateMultiple) GenerateMultipleVueTemplatesFunc(w http.ResponseWriter, r *http.Request, client interfaces.OpenAIClient, filesTextContent []models.VueFile) models.GenerateMultipleVueTemplateResponse {
	if s.execute != nil {
		println("EXECUTION")
		result := s.execute(w, r, client, filesTextContent)
		return result
	}
	return models.GenerateMultipleVueTemplateResponse{}
}

func (s *MockParseFiles) RequestParseFilesFunc(r *http.Request, w http.ResponseWriter) []*multipart.FileHeader {
	if s.execParseFiles != nil {
		return s.execParseFiles(r, w)
	}
	handler := helpers.RequestParseFiles{}
	return handler.RequestParseFilesFunc(r, w)
}
func (s *MockGetTextContent) GetTextContentFromFiles(files []interfaces.FileHeader, w http.ResponseWriter) ([]models.VueFile, error) {
	if s.execGetTextContent != nil {
		return s.execGetTextContent(files, w)
	}
	return helpers.GetTextContentFromFiles(files, w)
}

type TestCase struct {
	name       string
	mock       func(s *MockGenerateMultiple, d *MockGetTextContent, f *MockParseFiles)
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
			mock:       func(s *MockGenerateMultiple, d *MockGetTextContent, f *MockParseFiles) {},
			givenFiles: []File{},
			assertion: func(r *httptest.ResponseRecorder, t *testing.T) {
				if r.Code != http.StatusBadRequest {
					t.Error("Not a 400 error code")
				}
				if strings.TrimSpace(r.Body.String()) != "No files uploaded" {
					t.Errorf("Wrong message returned %s", r.Body.String())
				}
			},
		},
		{
			name: "Case 2: When request goes through, header should be set, correct status code returned, and contetn should be correct",
			mock: func(s *MockGenerateMultiple, d *MockGetTextContent, f *MockParseFiles) {
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
				if r.Header().Get("Content-Type") != "application/json" {
					t.Error("Doenst have application/json header")
				}
				if r.Code != http.StatusOK {
					t.Errorf("Doesnt have correct status code, expected: %v, is: %v", http.StatusOK, r.Code)
				}
			},
		},
		{
			name: "Case 3: Create error when json marshaling fails",
			mock: func(s *MockGenerateMultiple, d *MockGetTextContent, f *MockParseFiles) {
				s.execute = func(w http.ResponseWriter, r *http.Request, client interfaces.OpenAIClient, files []models.VueFile) models.GenerateMultipleVueTemplateResponse {
					// TODO: Create a response that will break json.Marshal()
					return models.GenerateMultipleVueTemplateResponse{models.GenerateSingleVueTemplateResponse{Content: ""}}
				}
			},
		},
		{
			name: "Case 4: Write json response data.",
			mock: func(s *MockGenerateMultiple, d *MockGetTextContent, f *MockParseFiles) {
				s.execute = func(w http.ResponseWriter, r *http.Request, client interfaces.OpenAIClient, files []models.VueFile) models.GenerateMultipleVueTemplateResponse {
					// TODO: Create a response that will break json.Marshal()
					return models.GenerateMultipleVueTemplateResponse{models.GenerateSingleVueTemplateResponse{Content: ""}}
				}
			},
		},
		// {
		// 	name: "If error during writing occurs create the correct error",
		// 	mock: func(s *MockGenerateMultiple) {
		// 		s.execute = func(w http.ResponseWriter, r *http.Request, client *openai.Client, filesTextContent []string) {
		// 			w.WriteHeader(http.StatusInternalServerError)
		// 		}
		// 	},
		// },
		{
			name: "Case 6: Error when extraction fails",
			mock: func(s *MockGenerateMultiple, d *MockGetTextContent, f *MockParseFiles) {
				d.execGetTextContent = func(files []interfaces.FileHeader, w http.ResponseWriter) ([]models.VueFile, error) {
					return nil, errors.New("Error occured")
				}
			},
			givenFiles: []File{{filename: "fileName1", content: "Diese"}, {filename: "fileName2", content: "Das ist ein Test"}},
			assertion: func(r *httptest.ResponseRecorder, t *testing.T) {
				if r.Code != http.StatusBadRequest {
					t.Error("Wrong code")
				}

				if r.Body.String() != "Error occured" {
					t.Error("Wrong error")
				}
			},
		},
	}

	for _, tc := range testCases {
		mockGenerateMultiple := MockGenerateMultiple{}
		mockParseFiles := MockParseFiles{}
		mockGetTextContent := MockGetTextContent{}

		tc.mock(&mockGenerateMultiple, &mockGetTextContent, &mockParseFiles)

		multipleFiles := MultipleFiles{
			GenerateMultipleVueTemplates: &mockGenerateMultiple,
			GetTextContentFromFiles:      &mockGetTextContent,
			RequestParseFiles:            &mockParseFiles,
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

	return req, nil
}
