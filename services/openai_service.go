package services

import (
	"context"
	"vue-converter-backend/interfaces"

	"github.com/sashabaranov/go-openai"
)

func GetVueChatCompletion(fileContent string) string {
	message := "This code is VueJS code that is not Version 3 and is not using the composition API. Please change this code that composition API is implemented. Just return the code, no explaining text. This is the code" + fileContent + "\n" + "Always make sure the whole component is returned including template, script and style. If you come across special properties prefixed with a '$' make sure to destructure it from the context parameter in the setup function."
	return message
}

func GetChatRequest(fileContent string) openai.ChatCompletionRequest {
	return openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo-16k",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "user",
				Content: GetVueChatCompletion(fileContent),
			},
		},
	}
}

type GenerateSingleTemplateResponseInterface interface {
	GenerateSingleTemplateResponse(request openai.ChatCompletionRequest, client interfaces.OpenAIClient) (openai.ChatCompletionResponse, error)
}

type GenerateSingleTemplateResponse struct{}

func (*GenerateSingleTemplateResponse) GenerateSingleTemplateResponseFunc(request openai.ChatCompletionRequest, client interfaces.OpenAIClient) (openai.ChatCompletionResponse, error) {
	resp, err := client.CreateChatCompletion(context.Background(), request)

	return resp, err
}
