package helpers

import "vue-converter-backend/interfaces"

func RemoveAuth0Prefix(auth0String string, compiler interfaces.RegexpCompile) (string, error) {
	regex, err := compiler.Compile(`^auth0\|`)
	if err != nil {
		return "", err
	}
	result := regex.ReplaceAllString(auth0String, "")
	return result, nil
}
