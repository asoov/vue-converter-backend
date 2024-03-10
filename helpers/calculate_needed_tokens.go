package helpers

import (
	"net/http"

	"github.com/pandodao/tokenizer-go"
)

func CalculateNeededTokens(input string, responseWriter http.ResponseWriter) int {
	tokensNeeded, err := tokenizer.CalToken(input)

	if err != nil {
		http.Error(responseWriter, "Error calculating needed tokens", http.StatusInternalServerError)
	}

	return tokensNeeded
}
