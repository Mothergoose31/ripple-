package config

import (
	"strings"
)

type Config struct {
	Mongo MongoConfig
	Path  PathConfig
}

type PathConfig struct {
	Secrets string
}

type MongoConfig struct {
	AuthEnabled        bool
	Host               string
	Port               int
	ConnectionInterval int
}

var App *Config

func GetConfig(configsetting string) {
	switch strings.ToUpper(configsetting) {
	case "DEVELOPMENT":
		App = DevelopmentConfig()
	case "TESTING":
		App = TestingConfig()
	}
}
