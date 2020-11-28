generate:
	@go generate

build: generate
	@go build -o gogreengrass main.go generated.go

install: build
	@go install

clean:
	@rm -rf generated.go