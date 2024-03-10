package main

import (
	"fmt"
	"net/http"
	"vue-converter-backend/env"
	"vue-converter-backend/routes"
)

func main() {
	env.InitializeEnvVars()
	routes.GenerateWebRoutes()

	fmt.Println("Server starting on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
