package handlers

import (
	"fmt"
	"net/http"
)

func GenerateMultiple(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Home Page!")
}
