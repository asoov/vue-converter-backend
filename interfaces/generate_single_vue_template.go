package interfaces

import (
	"net/http"
	"vue-converter-backend/models"
)

type GenerateSingleVueTemplate interface {
	GenerateSingleVueTemplateFunc(w http.ResponseWriter, r *http.Request, client OpenAIClient) models.GenerateSingleVueTemplateResponse
}
