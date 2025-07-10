run:
	go mod tidy
	go run cmd/server/main.go
daemon:
	CompileDaemon --build="go build -o server ./cmd/server" --command="./server"
