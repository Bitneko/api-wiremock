module github.com/Bitneko/api-wiremock

go 1.12

require (
	api-wiremock/configuration v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.7.2
	github.com/labstack/gommon v0.2.9 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.4.0 // indirect
	github.com/twinj/uuid v1.0.0
)

replace api-wiremock/configuration => ./configuration
