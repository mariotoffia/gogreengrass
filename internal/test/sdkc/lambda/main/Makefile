deps:
	@cp ../../../../../sdkc/libaws-greengrass-core-sdk-c.so .

build:
	@go build -o testlambda main.go

run:
	@go run main.go

clean:
	@rm -f testlambda libaws-greengrass-core-sdk-c.so