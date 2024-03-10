package models

type GenerateSingleVueTemplateRequest struct {
	FileName   string `json:"fileName"`
	Content    string `json:"content"`
	CustomerID string `json:"customerId"`
}
