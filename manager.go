package main

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

func (cm *ClientManager) ProcessClientListEvents {
	
}
