package helpers

import (
	"mime/multipart"
	"net/http"
)

func RequestParseFiles(r *http.Request, w http.ResponseWriter) []*multipart.FileHeader {
	parseErr := r.ParseMultipartForm(10 << 20)

	if parseErr != nil {
		http.Error(w, parseErr.Error(), http.StatusBadRequest)
		println(parseErr.Error())
		return nil
	}

	files := r.MultipartForm.File["files"]

	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return nil
	}

	return files
}
