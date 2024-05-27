package helpers

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"vue-converter-backend/interfaces"
	"vue-converter-backend/models"
)

// Constructor for the mock file.
func NewMockMultipartFile(content []byte) *interfaces.MockMultipartFile {
	return &interfaces.MockMultipartFile{
		Content: content,
		Reader:  bytes.NewReader(content),
	}
}

type MockFileOpener struct {
	OpenFunc func() (multipart.File, error)
	FileName string
}

func (m *MockFileOpener) Open() (multipart.File, error) {
	return m.OpenFunc()
}

func (m *MockFileOpener) Filename() string {
	return m.FileName
}
func TestGetTextContentFromFiles(t *testing.T) {
	type TestParameters struct {
		files []interfaces.FileHeader
		w     http.ResponseWriter
	}

	type TestExpectedResult struct {
		fileContents []models.VueFile
		err          error
	}

	type testCase struct {
		description string
		parameters  TestParameters
		expected    TestExpectedResult
	}

	testCases := []testCase{
		{
			description: "File",
			parameters: TestParameters{
				files: []interfaces.FileHeader{
					&MockFileOpener{
						OpenFunc: func() (multipart.File, error) { return NewMockMultipartFile([]byte("")), errors.New("error") },
						FileName: "Di",
					},
				},
				w: httptest.NewRecorder()},
			expected: TestExpectedResult{
				fileContents: nil,
				err:          errors.New("error"),
			},
		},
		{
			description: "Test case 1",
			parameters: TestParameters{
				files: []interfaces.FileHeader{
					&MockFileOpener{
						OpenFunc: func() (multipart.File, error) { return NewMockMultipartFile([]byte("hallo")), nil },
						FileName: "Diese",
					},
				},
				w: httptest.NewRecorder()},
			expected: TestExpectedResult{
				fileContents: []models.VueFile{{Name: "Diese", Content: "hallo"}},
				err:          nil,
			},
		},
	}

	// Iterate through test cases if needed
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			getTextContentFromFiles := GetTextContentFromFiles{}
			fileContents, err := getTextContentFromFiles.GetTextContentFromFilesFunc(tc.parameters.files, tc.parameters.w)

			if err != nil && err.Error() != tc.expected.err.Error() {
				t.Errorf("%s: Expected error to be %v but got %v", tc.description, tc.expected.err, err)
			} else if !reflect.DeepEqual(fileContents, tc.expected.fileContents) {
				t.Errorf("%s: Expected fileContents to be equal: %v, got: %v", tc.description, tc.expected.fileContents, fileContents)
			}
		})
	}
}
