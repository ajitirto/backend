run:
	go run ./cmd/api

build:
	go build -o bin/app ./cmd/api

docker:
	docker compose up --build

test:
	go test ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run

clean:
	rm -rf bin
