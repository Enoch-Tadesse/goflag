run:
	go mod tidy
	go run cmd/server/main.go
test:
	go test ./... -v
