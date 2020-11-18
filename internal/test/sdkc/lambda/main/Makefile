build:
	@mkdir -p ../_out/testlambda
	@go build -o ../_out/testlambda/testlambda main.go

run:
	@go run main.go

# https://docs.aws.amazon.com/greengrass/latest/developerguide/lambda-functions.html#lambda-executables
install:build
	@zip testlambda.zip testlambda
	aws lambda create-function \
	--function-name testlambda \
	--handler testlambda \
	--zip-file fileb://testlambda.zip \
	--runtime arn:aws:greengrass:::runtime/function/executable

	aws lambda publish-version \
	--function-name function-name testlambda

clean:
	@rm -rf ../_out/testlambda testlambda