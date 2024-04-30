// In the services package
package services

import (
	"net/http"
	"vue-converter-backend/interfaces"
	"vue-converter-backend/models"
)

type GenerateMultipleInterface interface {
	GenerateMultipleVueTemplatesFunc(w http.ResponseWriter, r *http.Request, client interfaces.OpenAIClient, files []models.VueFile) models.GenerateMultipleVueTemplateResponse
}
