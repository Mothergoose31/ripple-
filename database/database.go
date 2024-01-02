package database

import (
	"context"
	"fmt"
	"ripple/config"
	"time"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Client *mongo.Client

var connected bool

func Connect() {
	connected = false
	connectToClient()
	interval := SetInterval(connectToClient, 30000, false)
	for {
		if connected {
			interval <- true
			break
		}
	}

}

func connectToClient() {
	app := config.GetConfig()
	mongoURI := fmt.Sprintf("mongodb://%s:%d", app.Mongo.Host, app.Mongo.Port)

	// Log and create context
	log.Printf("Attempting connection on %s\n", mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set mongo client options
	clientOptions := options.Client().ApplyURI(mongoURI)
	if app.Mongo.AuthEnabled {
		clientOptions = clientOptions.SetAuth(getCredentials())
	}

	// Connect to mongo instance
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println(err)
		return
	}

	// Test mongo connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println(err)
		return
	}

	// Mark connection as successful and set client
	log.Println("Connection successful")
	connected = true
	Client = client
}
