package models

type GenerateMultipleVueTemplateResponse struct {
	FileName     string `json:"fileName"`
	Content      string `json:"content"`
	TokensNeeded int    `json:"tokensNeeded"`
	ErrorMessage string `json:"errorMessage"`
}
