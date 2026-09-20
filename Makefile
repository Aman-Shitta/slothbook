run:
	go run cmd/slotbook-api.go

test:
	echo "test not implemnted"

lint:
	echo "lint not implemnted"

fmt:
	go fmt

build:
	go build -o server ./cmd/
	./server
