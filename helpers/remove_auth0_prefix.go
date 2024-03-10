package helpers

import (
	"net/http"
	"regexp"
)

func RemoveAuth0Prefix(w http.ResponseWriter, auth0String string) string {
	regex, err := regexp.Compile(`^auth0\|`)
	if err != nil {
		http.Error(w, "Removing Auth0 prefix failed", http.StatusInternalServerError)
	}
	result := regex.ReplaceAllString(auth0String, "")
	return result
}
