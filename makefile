build:
	@go build -o gopad .

run: build
	@./gopad