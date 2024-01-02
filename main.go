package main

import (
	"ripple/config"
	"ripple/database"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {
	config.SetAppConfig("DEVELOPMENT")
	database.Connect()
}
