package config

import "strings"

//---------------------------------------------------------------

type Config struct {
	Mongo MongoConfig
	Path  PathConfig
}

//---------------------------------------------------------------

type PathConfig struct {
	Secrets string
}

//---------------------------------------------------------------

type MongoConfig struct {
	AuthEnabled        bool
	Host               string
	Port               int
	ConnectionInterval int
}

//---------------------------------------------------------------

var AppConfig *Config

func GetConfig() *Config {
	return AppConfig

}

//---------------------------------------------------------------

func SetAppConfig(configsetting string) {
	switch strings.ToUpper(configsetting) {
	case "DEVELOPMENT":
		AppConfig = DevelopmentConfig()
	case "TESTING":
		AppConfig = TestingConfig()
	}
}

//---------------------------------------------------------------

func SecretsPath() string {
	return AppConfig.Path.Secrets
}
