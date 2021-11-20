.PHONY: run tidy

run:
	go run cmd/main.go -env=local

tidy:
	go mod tidy
