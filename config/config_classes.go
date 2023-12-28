package config

func DevelopmentConfig() *Config {
	return &Config{
		Mongo: MongoConfig{
			AuthEnabled:        false,
			Host:               "localhost",
			Port:               27017,
			ConnectionInterval: 5000,
		},
		Path: PathConfig{
			Secrets: "../instance",
		},
	}
}

func TestingConfig() *Config {
	return &Config{
		Mongo: MongoConfig{
			AuthEnabled:        false,
			Host:               "localhost",
			Port:               27017,
			ConnectionInterval: 5000,
		},
		Path: PathConfig{
			Secrets: "../instance",
		},
	}
}
