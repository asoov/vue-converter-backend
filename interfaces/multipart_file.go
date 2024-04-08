package interfaces

import (
	"bytes"
	"mime/multipart"
)

// MockMultipartFile is our mock for multipart.File
type MockMultipartFile struct {
	Content []byte        // Content of the file
	Reader  *bytes.Reader // Reader to read Content, leveraging bytes.Reader
}

// Ensure that MockMultipartFile satisfies the multipart.File interface.
// Note: This verification step is causing the error because the method signatures
// are not fully satisfied as initially thought due to missing ReadAt directly.
var _ multipart.File = (*MockMultipartFile)(nil)

func (m *MockMultipartFile) Read(p []byte) (n int, err error) {
	return m.Reader.Read(p)
}

// Forward the Seek method to the bytes.Reader.
func (m *MockMultipartFile) Seek(offset int64, whence int) (int64, error) {
	return m.Reader.Seek(offset, whence)
}

// Implement the ReadAt method, forwarding to the bytes.Reader's ReadAt.
// This was missing in the initial implementation causing the interface satisfaction error.
func (m *MockMultipartFile) ReadAt(p []byte, off int64) (n int, err error) {
	return m.Reader.ReadAt(p, off)
}

// Close can be a no-op for the mock, as there are no open resources to dispose of.
func (m *MockMultipartFile) Close() error {
	return nil
}
