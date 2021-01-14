all: clean test

clean:
	@echo "Cleaning up..."
	@go mod tidy -v

test:
	@echo "Running tests..."
	@go test -v -cover ./...
