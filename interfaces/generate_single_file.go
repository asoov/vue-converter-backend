package interfaces

import (
	"net/http"

	"github.com/sashabaranov/go-openai"
)

type GenerateSingleFiles interface {
	GenerateSingle(
		w http.ResponseWriter,
		r *http.Request,
		client *openai.Client,
	)
}
