package routes

import (
	"net/http"
	"os"
	"vue-converter-backend/handlers"

	"github.com/sashabaranov/go-openai"
)

func generateSingleRoute(w http.ResponseWriter, r *http.Request) {
	client := openai.NewClient(os.Getenv("OAI_KEY"))
	singleFileHandler := handlers.GenerateSingleFile{}
	singleFileHandler.GenerateSingle(w, r, client)
}

func generateMultipleRoute(w http.ResponseWriter, r *http.Request) {
	client := openai.NewClient(os.Getenv("OAI_KEY"))
	multipleFiles := handlers.MultipleFiles{}
	multipleFiles.GenerateMultipleFiles(w, r, client)

}

func GenerateWebRoutes() {
	http.HandleFunc("/generate-web", generateSingleRoute)
	http.HandleFunc("generate-web/multiple", generateMultipleRoute)
	http.HandleFunc("generate-web/calulate-tokens", handlers.CalculateTokens)
}
