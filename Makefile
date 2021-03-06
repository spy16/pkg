all: clean fmt test

fmt:
	@echo "Formatting..."
	@goimports -l -w ./

clean:
	@echo "Cleaning up..."
	@go mod tidy -v

test:
	@echo "Running tests..."
	@go test -cover ./...

test-verbose:
	@echo "Running tests..."
	@go test -v -cover ./...

benchmark:
	@echo "Running benchmarks..."
	@go test -benchmem -run="none" -bench="Benchmark.*" -v ./...
