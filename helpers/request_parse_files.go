package helpers

import (
	"mime/multipart"
	"net/http"
)

type RequestParseFilesInterface interface {
	RequestParseFilesFunc(r *http.Request, w http.ResponseWriter) []*multipart.FileHeader
}

type RequestParseFiles struct{}

func (s *RequestParseFiles) RequestParseFilesFunc(r *http.Request, w http.ResponseWriter) []*multipart.FileHeader {
	parseErr := r.ParseMultipartForm(10 << 20)

	if parseErr != nil {
		http.Error(w, parseErr.Error(), http.StatusBadRequest)
		println(parseErr.Error())
		return nil
	}

	files := r.MultipartForm.File["files"]

	if len(files) == 0 {
		return nil
	}

	return files
}
