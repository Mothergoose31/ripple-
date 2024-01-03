package main

import (
	"context"

	"github.com/labstack/echo/v4"
)

type ClientList struct {
	EventType string
	Client    *Client
}

// if you want to be connected to multiple chatrooms then you need to have multiple clients
type ClientManager struct {
	Clients     map[string]map[string]*Client
	ClientLists chan *ClientList
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		Clients:     make(map[string]map[string]*Client),
		ClientLists: make(chan *ClientList),
	}
}

func (cm *ClientManager) ProcessClientListEvents(ctx context.Context) {
	for {
		select {
		case clientListEvent, ok := <-cm.ClientLists:
			if !ok {
				return
			}
			room := clientListEvent.Client.Chatroom

			// Initialize room if it doesn't exist
			if _, ok := cm.Clients[room]; !ok {
				cm.Clients[room] = make(map[string]*Client)
			}

			switch clientListEvent.EventType {
			case "ADD":
				// Check if client already exists in the room
				if _, exists := cm.Clients[room][clientListEvent.Client.ID]; exists {
					continue
				}
				cm.Clients[room][clientListEvent.Client.ID] = clientListEvent.Client

			case "REMOVE":
				// Remove the client from the room
				if _, exists := cm.Clients[room][clientListEvent.Client.ID]; exists {
					delete(cm.Clients[room], clientListEvent.Client.ID)
				}
			}

		case <-ctx.Done():
			return
		}
	}
}

// var upgrader = websocket.Upgrader{
//     ReadBufferSize:  1024,
//     WriteBufferSize: 1024,
// }

func (cm *ClientManager) Handle(c echo.Context, ctx context.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	newClient := NewClient(ws, cm, c.Param("chatroom"))

	// Sending the new client to the ClientListEventChannel
	cm.ClientLists <- &ClientList{
		EventType: "ADD",
		Client:    newClient,
	}

	// Start goroutines for reading and writing messages
	go newClient.ReadMessages(c)      // Ensure this method is defined in Client
	go newClient.WriteMessage(c, ctx) // Ensure this method is defined in Client

	return nil
}

func (cm *ClientManager) WriteMessage(msg string, chatroom string) error {
	// Check if the chatroom exists
	if clients, exists := cm.Clients[chatroom]; exists {
		// Iterate over all clients in the chatroom and send the message
		for _, client := range clients {
			client.MessageChan <- msg
		}
	}
	return nil
}
