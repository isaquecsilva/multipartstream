package multipartstream

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"testing"
)

const defaultErrFormat = "expected = %v, actual = %v"

func assertPanics(t *testing.T, f func()) {
	defer func() {
		if err := recover(); err == nil {
			t.FailNow()
		}
	}()

	f()
}

func Test_generateBoundary(t *testing.T) {
	ms := NewMultipartStream()
	re := regexp.MustCompile(`WebkitFormBoundary[0-9a-f]{20}`)

	boundary := ms.generateBoundary()
	if result := re.MatchString(boundary); !result {
		t.Errorf(defaultErrFormat, true, result)
	}
}

func TestBoundary(t *testing.T) {
	re := regexp.MustCompile(`WebkitFormBoundary[0-9a-f]{20}`)

	ms := NewMultipartStream()

	if ms.Boundary() == "" {
		t.Errorf(defaultErrFormat, "not-empty", ms.Boundary())
	}


	match := re.MatchString(ms.Boundary())

	if !match {
		t.Errorf(defaultErrFormat, true, match)
	}
}

func TestContentType(t *testing.T) {
	ms := NewMultipartStream()
	contentType := ms.ContentType()

	const expected = "multipart/form-data;"
	if !strings.Contains(contentType, expected) {
		t.Errorf(defaultErrFormat, expected, contentType)
	}
}

func TestFormField(t *testing.T) {
	ms := NewMultipartStream()

	assertPanics(t, func() {
		ms.FormField("some-field", nil)
	})

	ms.FormField("some-field", strings.NewReader("ok"))
	var expected = fmt.Sprintf("--%s\nContent-Disposition: form-data; name=\"some-field\"\n\nok\n",
		ms.Boundary(),
	)

	buf, _ := io.ReadAll(ms)
	if !bytes.Equal(buf, []byte(expected)) {
		t.Errorf(defaultErrFormat, expected, string(buf))
	}
}

func TestFormFile(t *testing.T) {
	ms := NewMultipartStream()

	// panics cause of nil reader
	assertPanics(t, func() {
		ms.FormFile("file", "example.txt", nil)
	})

	// error -> could not determine mimetype
	err := ms.FormFile("file", "example", strings.NewReader("some file contents..."))
	if err == nil {
		t.Errorf(defaultErrFormat, "not-nil", err)
	}

	// asserting error message: not extension
	var expected = "no extension in filename"
	if err.Error() != expected {
		t.Errorf(defaultErrFormat, expected, err.Error())
	}

	err = ms.FormFile("file", "example.asdioieor", strings.NewReader("some file contents"))

	if err == nil {
		t.Errorf(defaultErrFormat, "not-nil", err)
	}

	// asserting error message: failure parsing mimetype
	expected = "could not parse mimetype"
	if err.Error() != expected {
		t.Errorf(defaultErrFormat, expected, err.Error())
	}

	/// success
	err = ms.FormFile("file", "example.txt", strings.NewReader("some file contents..."))
	if err != nil {
		t.Errorf(defaultErrFormat, nil, err)
	}
}

func TestDone(t *testing.T) {
	ms := NewMultipartStream()
	ms.FormFile("file", "hello.txt", strings.NewReader("world!"))
	ms.Done()

	expected := []byte(ms.Boundary() + "--")

	buf, _ := io.ReadAll(ms)
	if !bytes.Contains(buf, expected) {
		t.Errorf(defaultErrFormat, string(expected), string(buf))
	}
}