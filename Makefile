test:
	@echo "Running tests..."
	@go clean -testcache
	@echo "Starting tests..."
	@START_TIME=$$(gdate +%s%3N); \
	GOEXPERIMENT=synctest go test -v ./...; \
	END_TIME=$$(gdate +%s%3N); \
	ELAPSED_MS=$$((END_TIME - START_TIME)); \
	echo "Tests completed in $$ELAPSED_MS ms."