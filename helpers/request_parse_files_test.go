package helpers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCase struct {
	description string
	testBody    func()
}

func TestRequestParseFiles(t *testing.T) {
	testCases := []TestCase{
		{
			description: "Case 1: Parsing Error",
			testBody: func() {
				request := httptest.NewRequest("POST", "/", nil)
				request.Header.Set("Content-Type", "multipart/form-data")

				w := httptest.NewRecorder()

				requestParseFiles := RequestParseFiles{}
				requestParseFiles.RequestParseFilesFunc(request, w)

				if w.Result().StatusCode != http.StatusBadRequest {
					t.Errorf("Theres an error %v", w.Result().StatusCode)
				}
				if w.Result().Status != "400 Bad Request" {
					t.Errorf("Failed %v", w.Result().Status)
				}
			},
		},
		{
			description: "Case 2: Parsing worked out",
			testBody: func() {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)
				file, _ := writer.CreateFormFile("files", "file.txt")

				file.Write([]byte("file content"))
				writer.Close()

				request := httptest.NewRequest("POST", "/", body)
				request.Header.Set("Content-Type", writer.FormDataContentType())

				w := httptest.NewRecorder()

				requestParseFiles := RequestParseFiles{}
				result := requestParseFiles.RequestParseFilesFunc(request, w)

				if w.Result().StatusCode != 200 {
					t.Errorf("Status code %d", w.Result().StatusCode)
				}
				if result == nil {
					t.Error("Parsing file did not work out!", result)
				}

				firstResultElement := result[0]
				firstResultFileName := firstResultElement.Filename

				if firstResultFileName != "file.txt" {
					t.Error("Filename is not correct", firstResultFileName)
				}

				firstResultFileOpen, openErr := firstResultElement.Open()
				if openErr != nil {
					t.Error("Failed to open file", openErr)
				}
				firstResultFileContent, readErr := io.ReadAll(firstResultFileOpen)
				if readErr != nil {
					t.Error("Failed to read file", readErr)
				}

				if string(firstResultFileContent) != "file content" {
					t.Error("File content is not correct", string(firstResultFileContent))
				}

			},
		},
		{
			description: "Case 3: No files uploaded",
			testBody: func() {
				body := new(bytes.Buffer)

				writer := multipart.NewWriter(body)
				writer.Close()

				request := httptest.NewRequest("POST", "/", body)
				request.Header.Set("Content-Type", writer.FormDataContentType())

				w := httptest.NewRecorder()

				requestParseFiles := RequestParseFiles{}
				result := requestParseFiles.RequestParseFilesFunc(request, w)

				if w.Result().StatusCode != http.StatusBadRequest {
					t.Errorf("Status code %d", w.Result().StatusCode)
				}

				if result != nil {
					t.Error("Parsing file did not work out!", result)
				}

			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			tc.testBody()
		})
	}

}
