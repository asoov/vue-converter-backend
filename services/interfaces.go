// services/interfaces.go

package services

import (
	"net/http"
	"vue-converter-backend/models"

	"github.com/sashabaranov/go-openai"
)

type ServicesInterface interface {
	GenerateSingleVueTemplate(w http.ResponseWriter, r *http.Request, client *openai.Client) models.GenerateSingleVueTemplateResponse
}
