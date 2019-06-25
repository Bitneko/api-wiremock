export GO111MODULE:=on

dev:
	$(eval export ENVIRONMENT=development)
	@go run ./...

clean:
	@rm -rf mocks/msf/.data/__files/*
	@rm -rf mocks/msf/.data/mappings/*