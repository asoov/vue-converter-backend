package helpers

import (
	"errors"
	"testing"
	"vue-converter-backend/models"
)

type MockTokenizerCalToken struct {
	replacement func(str string) (int, error)
}

func (f *MockTokenizerCalToken) CalToken(str string) (int, error) {
	if f.replacement != nil {
		return f.replacement(str)
	}
	return 6, nil
}

func TestCalculateNeededTokens(t *testing.T) {
	t.Run("It should return the number of tokens needed to generate the templates", func(t *testing.T) {

		var callCount int = 0
		calculateNeededTokens := CalculateNeededTokens{Tokenizer: &MockTokenizerCalToken{replacement: func(str string) (int, error) {
			defer func() { callCount++ }()
			if callCount == 0 {
				return 20, nil
			} else {
				return 6, nil
			}
		}}}

		vueFiles := []models.VueFile{{Name: "Element 1", Content: "This is the content 1"}, {Name: "Element 2", Content: "This is the content 2"}}

		result, _ := calculateNeededTokens.CalculateNeededTokensFunc(vueFiles)

		expectedTokenAmount := 26
		if result != expectedTokenAmount {
			t.Errorf("Expected %v, got %v", expectedTokenAmount, result)
		}
	})

	t.Run("It should return an error if the tokenizer returns an error", func(t *testing.T) {
		CalculateNeededTokens := CalculateNeededTokens{Tokenizer: &MockTokenizerCalToken{replacement: func(str string) (int, error) {
			return 0, errors.New("Tokenizer error")
		}}}

		vueFiles := []models.VueFile{{Name: "Element 1", Content: "This is the content 1"}, {Name: "Element 2", Content: "This is the content 2"}}
		_, err := CalculateNeededTokens.CalculateNeededTokensFunc(vueFiles)

		if err == nil {
			t.Error("Expected an error, got nil")
		}
	})

}
