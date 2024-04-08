package interfaces

import "regexp"

type RegexpCompile interface {
	Compile(str string) (*regexp.Regexp, error)
}
