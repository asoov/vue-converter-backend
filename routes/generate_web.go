package routes

import (
	"net/http"
	"vue-converter-backend/handlers"
)

func GenerateWebRoutes() {
	http.HandleFunc("/generate-web", handlers.GenerateSingle)
	http.HandleFunc("generate-web/multiple", handlers.GenerateMultiple)
	http.HandleFunc("generate-web/calulate-tokens", handlers.CalculateTokens)
}
