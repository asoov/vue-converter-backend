package services

import (
	"net/http"
	"vue-converter-backend/interfaces"
	"vue-converter-backend/models"
)

type GenerateMultipleInterface struct {
	Services ServicesInterface
}

func (s *GenerateMultipleInterface) GenerateMultipleVueTemplates(w http.ResponseWriter, r *http.Request, client interfaces.OpenAIClient, filesTextContent []string) models.GenerateMultipleVueTemplateResponse {
	return models.GenerateMultipleVueTemplateResponse{}
}
