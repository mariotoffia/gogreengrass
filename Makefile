generate:
	@go generate

build: generate
	@go build -o gogreengrass main.go generated.go

install: build
	@go install

clean:
	@rm -f gogreengrass generated.go glue.go glue.py
	@rm -rf greengrasssdk