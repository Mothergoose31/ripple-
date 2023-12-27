package main

type ClientList struct {
	EventType string
	Client    *Client
}

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
