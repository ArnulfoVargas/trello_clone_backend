run:
	@go run ./cmd/main.go
tidy:
	@go mod tidy
build:
	@templ generate
	@go build cmd/main.go
	@./main.exe