package handlers

import (
	"encoding/json"
	"net/http"
	"vue-converter-backend/services"
)

func GenerateSingle(w http.ResponseWriter, r *http.Request) {

	response := services.GenerateSingleVueTemplate(w, r)
	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(response)

	if err != nil {
		// If an error occurs during JSON marshaling, send an error to the client
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, writeErr := w.Write(jsonData)
	if writeErr != nil {
		// If there is an error while writing, log it, or handle it as necessary
		http.Error(w, "Failed to write JSON to response", http.StatusInternalServerError)
	}
}