run:
	go run ./cmd/demo

test:
	go test ./... -v -race

fmt:
	go fmt ./...
