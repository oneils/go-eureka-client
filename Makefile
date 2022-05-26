.PHONY:
.SILENT:

test:
	go test --short -coverprofile=cover.out ./...
	make test.coverage

test.coverage:
	go tool cover -func=cover.out

lint:
	golangci-lint run

.DEFAULT_GOAL := build%