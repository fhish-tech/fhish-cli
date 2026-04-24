BINARY_NAME=fhish

build:
	go build -o bin/$(BINARY_NAME) main.go

install:
	go install

clean:
	rm -rf bin/
	go clean

run:
	go run main.go

test:
	go test ./...

.PHONY: build install clean run test
