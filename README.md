# MultipartStream

Multipartstream offers multipart/form-data body requests as stream object, avoiding overheads on upload processes by larger files.

**Its currently under development, so there's no garrantees of stability or backwards compatibility.**

# Installation
```go
go get github.com/isaquecsilva/multipartstream
```

##

## Example Usage:

```go
package main

import (
    "os"
    "strings"
    "net/http"

    ms "github.com/isaquecsilva/multipartstream"
)

func main() {
    // new multipartStream instance
    form := ms.NewMultipartStream()

    // adding a new form field called 'sender'
    form.FormField("sender", strings.NewReader("Herold"))
    
    // loading file to upload
    example, err := os.Open("example.txt")
    if err != nil {
        panic(err)
    }
    defer example.Close()

    // adding our example.txt file, and checking for mimetypes errors.
    if err = form.FormFile("file", "example.txt", example); err != nil {
    	panic(err)
    }

    // adding trailing boundary signaling multipart/form-data body end.
    form.Done()

    // our http request    
    req, _ := http.NewRequest(http.MethodPost, endpoint, form)

    // sending request
    http.DefaultClient.Do(req)
}
```

# Inspiration
It was inspired on famous Go package <a target="_blank" href="github.com/technoweenie/multipartstreamer"><strong>technoweenie/multipartstreamer</strong></a>

# License
- Terms of use for Multipartstream resides under MIT license.


