all: build

generate: clean
	@gogreengrass -d

build: generate
	@go build -o main.so -buildmode=c-shared main.go glue.go

clean:
	@rm -f main.so main.h glue.go glue.py > /dev/null 2>&1
	@rm -rf __pycache__ greengrasssdk> /dev/null 2>&1