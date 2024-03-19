package services

import (
	"net/http"
	"vue-converter-backend/interfaces"
	"vue-converter-backend/models"
)

type GenerateMultipleVueTemplateFunc func(w http.ResponseWriter, r *http.Request, client interfaces.OpenAIClient) models.GenerateMultipleVueTemplateResponse

func GenerateMultipleVueTemplates(w http.ResponseWriter, r *http.Request, client interfaces.OpenAIClient) models.GenerateMultipleVueTemplateResponse {
	return models.GenerateMultipleVueTemplateResponse{}
}
