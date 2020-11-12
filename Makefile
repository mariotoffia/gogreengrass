all: run

.SUFFIXES: .go .so
.go.so:
	@go build -o $@ -buildmode=c-shared $<

build: main.so

run: build
	@python3 -m pycode

clean:
	@rm -f main.so main.h > /dev/null 2>&1
	@rm -rf __pycache__ > /dev/null 2>&1