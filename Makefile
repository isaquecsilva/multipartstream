test:
	@go test ./... -v

cover:
	@go test ./... -cover

example:
	@go run examples/simple_file_upload.go