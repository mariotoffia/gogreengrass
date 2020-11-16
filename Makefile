generate:
	@go generate

build: generate
	@go build main.go generated.go

run:
	@go run main.go generated.go
clean:
	@rm -f gogreengrass generated.go glue.go glue.py