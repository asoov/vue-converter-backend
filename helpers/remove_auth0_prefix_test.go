package helpers

import (
	"regexp"
	"testing"
)

type MockRegexp struct{}

func (m *MockRegexp) Compile(str string) (*regexp.Regexp, error) {
	return regexp.MustCompile(`^auth0\|`), nil
}

func TestRemoveAuth0Prefix(t *testing.T) {
	mockCompiler := &MockRegexp{}
	result, err := RemoveAuth0Prefix("hallo", mockCompiler)

	if err != nil {
		t.Error("faileddd", result)
	}
}
