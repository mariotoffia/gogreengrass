generate:
	@go generate

build: generate
	@go build

clean:
	@rm -f gogreengrass generated.go