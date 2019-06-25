package configuration

import (
	"fmt"
	"reflect"

	"github.com/spf13/viper"
)

// IConfig describe Config's type
type IConfig struct {
	APIWiremock string
	APITarget   string
	ProxyURL    string
	Environment string
}

// Config describe environment setting
var Config = IConfig{}

// InitConfig read in environment variables in env.yaml or system
func InitConfig() {
	viper.AutomaticEnv()
	getEnvironment()

	if Config.Environment == "development" {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if err != nil { // Handle errors reading the config file
			panic(fmt.Errorf("fatal error config file: %s", err))
		}
	}

	constructConfig()

	initLogger()
}

func getEnvironment() {
	Config.Environment = viper.GetString("ENVIRONMENT")

	if Config.Environment == "" {
		Config.Environment = "Environment Not Specified"
	}
}

func constructConfig() {
	Config.APIWiremock = viper.GetString("API_WIREMOCK")
	Config.APITarget = viper.GetString("API_TARGET")
	Config.ProxyURL = viper.GetString("PROXY_URL")

	if Config.Environment == "development" {
		fmt.Println("====== Environment Variables ======")
		values := reflect.ValueOf(&Config).Elem()
		keys := values.Type()

		for i := 0; i < values.NumField(); i++ {
			fmt.Println(keys.Field(i).Name, "=", values.Field(i))
		}
	}
}
