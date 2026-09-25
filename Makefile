run:
	go run cmd/slotbook-api.go

test:
	go test ./internal/test/

lint:
	echo "lint not implemnted"

fmt:
	go fmt

build:
	go mod tidy
	go build -o server ./cmd/
	./server
