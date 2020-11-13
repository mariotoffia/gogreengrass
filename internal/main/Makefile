all: run

build: 
	@go build -o main.so -buildmode=c-shared main.go

run: build
	@python3 -m pycode

clean:
	@rm -f main.so main.h > /dev/null 2>&1
	@rm -rf __pycache__ > /dev/null 2>&1