package routes

import (
	"net/http"
	"os"
	"vue-converter-backend/handlers"
	"vue-converter-backend/services"

	"github.com/sashabaranov/go-openai"
)

var generateSingleRoute = func(w http.ResponseWriter, r *http.Request) {
	client := openai.NewClient(os.Getenv("OAI_KEY"))
	var generateFunc services.GenerateSingleVueTemplateFunc = services.GenerateSingleVueTemplate
	handlers.GenerateSingle(w, r, client, generateFunc)
}

func generateMultipleRoute(w http.ResponseWriter, r *http.Request) {
	client := openai.NewClient(os.Getenv("OAI_KEY"))
	var generateFunc services.GenerateMultipleVueTemplateFunc = services.GenerateMultipleVueTemplates
	handlers.GenerateMultiple(w, r, client, generateFunc)
}

func GenerateWebRoutes() {
	http.HandleFunc("/generate-web", generateSingleRoute)
	http.HandleFunc("generate-web/multiple", generateMultipleRoute)
	http.HandleFunc("generate-web/calulate-tokens", handlers.CalculateTokens)
}
