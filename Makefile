build:
	go build -o bin/calc cmd/calc/main.go

run:
	go run cmd/srv/*.go

deps:
	@echo "downloading dependencies"
	go mod tidy
	go mod vendor
	go mod download
