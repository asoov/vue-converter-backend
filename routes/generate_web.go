package routes

import (
	"net/http"
	"os"
	"vue-converter-backend/handlers"

	"github.com/sashabaranov/go-openai"
)

var generateSingleRoute = func(w http.ResponseWriter, r *http.Request) {
	client := openai.NewClient(os.Getenv("OAI_KEY"))
	handlers.GenerateSingle(w, r, client)
}

func GenerateWebRoutes() {
	http.HandleFunc("/generate-web", generateSingleRoute)
	http.HandleFunc("generate-web/multiple", handlers.GenerateMultiple)
	http.HandleFunc("generate-web/calulate-tokens", handlers.CalculateTokens)
}
