package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"

	ms "github.com/isaquecsilva/multipartstream"
)

func main() {
	srv := testServer()
	defer srv.Close()

	wd, _ := os.Getwd()
	example, err := os.Open(filepath.Join(wd, "examples/example.txt"))
	if err != nil {
		log.Fatal(err)
	}
	defer example.Close()

	form := ms.NewMultipartStream()
	form.FormField("whosent", strings.NewReader("oliverwho"))
	form.FormFile("file", filepath.Base(example.Name()), example)
	form.Done()

	req, _ := http.NewRequest(http.MethodPost, srv.URL, form)
	req.Header.Add("Content-Type", form.ContentType())

	http.DefaultClient.Do(req)
}

func handler() http.Handler {
	const maxBytes int64 = 1000_000

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if err := r.ParseMultipartForm(maxBytes); err != nil {
			log.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}


		sender := r.MultipartForm.Value["whosent"][0]
		log.Printf("Sender = %s\n", sender)

		file, header, err := r.FormFile("file")

		if err != nil {
			log.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		defer file.Close()

		log.Printf("Filename = %s, Size = %d\n", header.Filename, header.Size)
		io.Copy(os.Stdout, file)
		w.WriteHeader(http.StatusOK)
	})
}

func testServer() *httptest.Server {
	return httptest.NewServer(handler())
}
