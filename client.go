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

//=====================================

func NewClient(conn *websocket.Conn, manager *ClientManager, chatroom string) *Client {
	return &Client{
		Connection:  conn,
		ID:          conn.RemoteAddr().String(),
		Chatroom:    chatroom,
		Manager:     manager,
		MessageChan: make(chan string, 100),
	}
}

//=====================================

func (c *Client) ReadMessages(ctx echo.Context) {
	if err := c.setReadDeadline(); err != nil {
		log.Println("SetReadDeadline error:", err)
		return
	}

	c.setupPongHandler(ctx)

	defer c.cleanupClient()

	for {
		if err := c.processIncomingMessage(); err != nil {
			log.Println("ReadMessage error:", err)
			return
		}
	}
}

func (c *Client) setReadDeadline() error {
	// Assuming writeWait is a properly defined duration
	return c.Connection.SetReadDeadline(time.Now().Add(pingWait))
}

// =====================================
func (c *Client) setupPongHandler(ctx echo.Context) {
	if err := c.setReadDeadline(); err != nil {
		ctx.Logger().Error(err)
		return
	}
	c.Connection.SetPongHandler(func(data string) error {
		if err := c.setReadDeadline(); err != nil {
			ctx.Logger().Error(err)
			return err
		}
		log.Println("Pong received from client:", c.ID)
		return nil
	})
	defer func() {
		c.Connection.Close()
		c.Manager.ClientLists <- &ClientList{
			Client:    c,
			EventType: "REMOVE",
		}
	}()
	for {
		_, msg, err := c.Connection.ReadMessage()
		if err != nil {
			ctx.Logger().Error(err)
			return
		}
		log.Printf("%s\n", msg)
		if err := c.Manager.WriteMessage(string(msg), c.Chatroom); err != nil {
			ctx.Logger().Error(err)
			return
		}
	}
}

//=====================================

func (c *Client) cleanupClient() {
	c.Connection.Close()
	c.Manager.ClientLists <- &ClientList{
		Client:    c,
		EventType: "REMOVE",
	}
}

//=====================================

func (c *Client) processIncomingMessage() error {
	_, msg, err := c.Connection.ReadMessage()
	if err != nil {
		return err
	}
	c.Manager.WriteMessage(string(msg), c.Chatroom)
	return nil
}

//=====================================

func (cm *ClientManager) WriteMessage(msg string, chatroom string) error {
	if clients, ok := cm.Clients[chatroom]; ok {
		for _, client := range clients {
			select {
			case client.MessageChan <- msg:
			// message sent successfully

			default:
				// handle error
				log.Printf("Failed to send message to client [%s] in chatroom [%s]", client.ID, chatroom)
			}
		}
	}
	return nil
}
