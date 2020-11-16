all: invoke shadow

invoke:
	@go build -o main.so -buildmode=c-shared invoke.go glue.go
	@python3 invoke.py

shadow:
	@go build -o main.so -buildmode=c-shared shadow.go glue.go
	@python3 shadow.py

clean:
	@rm -f main.so main.h > /dev/null 2>&1
	@rm -rf __pycache__ > /dev/null 2>&1