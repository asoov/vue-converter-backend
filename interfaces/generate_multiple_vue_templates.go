package interfaces

import (
	"net/http"
	"vue-converter-backend/models"
)

type GenerateMultipleVueTemplates interface {
	GenerateMultipleVueTemplatesFunc(w http.ResponseWriter, r *http.Request, client OpenAIClient, files []models.VueFile) (models.GenerateMultipleVueTemplateResponse, int, error)
}
