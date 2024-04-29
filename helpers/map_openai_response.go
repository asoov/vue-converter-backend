package helpers

import (
	"vue-converter-backend/models"

	"github.com/sashabaranov/go-openai"
)

func MapOpenAiResponse(fileName string, response openai.ChatCompletionResponse, responseError error) models.GenerateSingleVueTemplateResponse {
	if responseError != nil {
		return models.GenerateSingleVueTemplateResponse{
			ErrorMessage: responseError.Error(),
		}
	}

	return models.GenerateSingleVueTemplateResponse{
		FileName:     fileName,
		Content:      response.Choices[0].Message.Content,
		TokensNeeded: response.Usage.TotalTokens,
	}
}
