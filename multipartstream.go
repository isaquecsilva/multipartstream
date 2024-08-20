package multipartstream

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"
)

const rawBoundarySize = 10

type MultipartStream struct {
	r        io.Reader
	boundary string
}

// NewMultipartStream creates a new instance os *MultipartStream object.
func NewMultipartStream() *MultipartStream {
	ms := new(MultipartStream)
	ms.r = new(bytes.Buffer)
	ms.boundary = ms.generateBoundary()
	return ms
}

func (ms *MultipartStream) generateBoundary() string {
	buf := make([]byte, rawBoundarySize)
	io.ReadFull(rand.Reader, buf)
	return fmt.Sprintf("WebkitFormBoundary%x", buf)
}

// Boundary returns MultipartStream multipart/form-data boundary value.
func (ms *MultipartStream) Boundary() string {
	return ms.boundary
}

// ContentType returns an HTTP Content-Type header value according to multipart/form-data with the generated boundary by NewMultipartStream.
func (ms *MultipartStream) ContentType() string {
	return fmt.Sprintf("multipart/form-data; boundary=%s", ms.Boundary())
}

// Read implements io.Reader interface.
func (ms *MultipartStream) Read(buf []byte) (int, error) {
	return ms.r.Read(buf)
}

// FormField adds a new form field to current multipart/form-data body. It panics if `r` is nil.
func (ms *MultipartStream) FormField(fieldName string, r io.Reader) {

	if r == nil {
		panic("FormField `r`, value of type io.Reader, is nil")
	}

	buffer := new(bytes.Buffer)
	buffer.WriteString("--" + ms.Boundary() + "\n")
	buffer.WriteString(fmt.Sprintf("Content-Disposition: form-data; name=\"%s\"\n\n", fieldName))

	io.Copy(buffer, r)
	buffer.WriteString("\n")

	if ms.r == nil {
		ms.r = buffer
		return
	}

	ms.r = io.MultiReader(ms.r, buffer)
}

// FormFile adds a new form file field in current multipart/form-data body. It panics if `r` is nil. The returned error, if there is one, will be related with mime type parsing failures.
func (ms *MultipartStream) FormFile(fieldName, filename string, r io.Reader) error {
	if r == nil {
		panic("FormFile `r`, value of type io.Reader, is nil")
	}

	buffer := new(bytes.Buffer)
	buffer.WriteString("--"+ms.Boundary() + "\n")
	buffer.WriteString(fmt.Sprintf("Content-Disposition: form-data; name=\"file\"; filename=\"%s\"\n", filename))

	mimeType, err := ms.getFileExtension(filename)
	if err != nil {
		return err
	}

	// setting headers and file reader into ms.r
	buffer.WriteString(fmt.Sprintf("Content-Type: %s\n\n", mimeType))
	ms.r = io.MultiReader(ms.r, buffer, r, strings.NewReader("\n"))
	return nil
}

func (ms *MultipartStream) getFileExtension(filename string) (string, error) {
	ext := filepath.Ext(filename)
	if ext == "" {
		return "", errors.New("no extension in filename")
	}

	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "", errors.New("could not parse mimetype")
	}

	return mimeType, nil
}

// Done inserts trailling boundary signaling the end of the multipart/form-data body request.
func (ms *MultipartStream) Done() {
	ms.r = io.MultiReader(ms.r, strings.NewReader("--"+ms.Boundary()+"--"))
}
