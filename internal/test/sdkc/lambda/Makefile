generate:
	@cd main && make generate && cd ..

build:
	@cd main && make build & cd ..
	@cd cdk && npm run build && cd ..

deploy:
	@cd cdk && npm run deploy && cd ..

destroy:
	@cd cdk && npm run destroy && cd ..
	
clean:
	@cd main && make clean