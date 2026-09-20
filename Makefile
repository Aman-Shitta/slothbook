run:
	go run cmd/main.go

test:
	echo "test not implemnted"

lint:
	echo "lint not implemnted"

fmt:
	go fmt

build:
	go build -o server ./cmd/main.go
	./server
