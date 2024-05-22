package interfaces

import "mime/multipart"

type FileHeader interface {
	Open() (multipart.File, error)
	Filename() string
}
