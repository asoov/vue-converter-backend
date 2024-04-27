package interfaces

import (
	"net/http"

	"github.com/sashabaranov/go-openai"
)

type GenerateMultipleFiles interface {
	GenerateMultipleFiles(w http.ResponseWriter, r *http.Request, client *openai.Client)
}
