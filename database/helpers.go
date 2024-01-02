package database

import (
	"os"
	"path/filepath"
	"ripple/config"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ==================================================

func readCredentials(file string) string {
	path := filepath.Join(config.AppConfig.Path.Secrets, file)

	credentials, err := os.ReadFile(path)
	if err != nil {
		log.Error(err)

	}
	return string(credentials)
}

// ==================================================

func getCredentials() options.Credential {
	return options.Credential{
		AuthMechanism: "SCRAM-SHA-256",
		AuthSource:    readCredentials("database"),
		Username:      readCredentials("username"),
		Password:      readCredentials("password"),
	}
}

// ==================================================

func SetInterval(someFunc func(), milliseconds int, async bool) chan bool {
	interval := time.Duration(milliseconds) * time.Millisecond
	ticker := time.NewTicker(interval)
	clear := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				if async {
					go someFunc()
				} else {
					someFunc()
				}
			case <-clear:
				ticker.Stop()
				return
			}
		}
	}()

	return clear
}
