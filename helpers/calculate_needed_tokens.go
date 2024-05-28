package helpers

import (
	"vue-converter-backend/interfaces"
	"vue-converter-backend/models"
)

type CalculateNeededTokensInterface interface {
	CalculateNeededTokensFunc(inputFiles []models.VueFile) (*int, error)
}

type CalculateNeededTokens struct {
	Tokenizer interfaces.Tokenizer
}

func (s *CalculateNeededTokens) CalculateNeededTokensFunc(textContents []models.VueFile) (int, error) {

	tokensNeeded := 0

	for _, file := range textContents {
		tokensNeededForFile, err := s.Tokenizer.CalToken(file.Content)

		if err != nil {
			return 0, err
		}

		tokensNeeded += tokensNeededForFile
	}

	return tokensNeeded, nil
}
