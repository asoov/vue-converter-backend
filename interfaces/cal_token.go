package interfaces

import "github.com/pandodao/tokenizer-go"

type Tokenizer interface {
	CalToken(input string) (int, error)
}

type TokenizerCalToken struct{}

func (f *TokenizerCalToken) CalToken(input string) (int, error) {
	tokensNeeded, err := tokenizer.CalToken(input)
	return tokensNeeded, err
}
