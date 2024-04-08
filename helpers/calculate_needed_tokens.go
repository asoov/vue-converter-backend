package helpers

func CalculateNeededTokens(input string, calculateTokens func(string) (int, error)) (int, error) {
	tokensNeeded, err := calculateTokens(input)

	if err != nil {
		return 0, err
	}

	return tokensNeeded, nil
}
