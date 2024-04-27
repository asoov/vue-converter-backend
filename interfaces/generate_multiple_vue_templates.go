package interfaces

import (
	"net/http"
	"vue-converter-backend/models"
)

type GenerateMultipleVueTemplates interface {
	GenerateMultipleVueTemplates(w http.ResponseWriter, r *http.Request, client OpenAIClient, filesTextContent []string) models.GenerateMultipleVueTemplateResponse
}
