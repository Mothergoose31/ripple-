package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Client struct {
	Connection  *websocket.Conn
	ID          string
	Chatroom    string
	Manager     *ClientManager
	MessageChan chan string
}

var (
	pingWait  = 60 * time.Second
	writeWait = 10 * time.Second
)

func NewClient(conn *websocket.Conn, manager *ClientManager, chatroom string) *Client {
	return &Client{
		Connection:  conn,
		ID:          conn.RemoteAddr().String(),
		Chatroom:    chatroom,
		Manager:     manager,
		MessageChan: make(chan string, 100),
	}
}

func (c *Client) ReadMessages(ctx echo.Context) {
	if err := c.Connection.SetReadDeadline(time.Now().Add(pingWait)); err != nil {
		log.Println("SetReadDeadline error:", err)
		return
	}
	//  writewait logic needs to be added
}

func (cm *ClientManager) WriteMessages(msg string, chatroom string) error {
	if clients, ok := cm.Clients[chatroom]; ok {
		for _, client := range clients {
			client.MessageChan <- msg
		}
	}
	return nil
}
