export GO111MODULE:=on

start:
	@docker-compose up

dev:
	@docker-compose start api_wiremock
	$(eval export ENVIRONMENT=development)
	@go run ./cmd/api-wiremock

clean:
	@rm -rf mocks/msf/.data/__files/*
	@rm -rf mocks/msf/.data/mappings/*