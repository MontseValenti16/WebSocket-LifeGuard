package service

import (
	"log"
	"sync"
	"WS/websoquet-server/internal/domain"

	"github.com/gorilla/websocket"
)

// WebsoquetService gestiona las conexiones de clientes identificados por accountID.
type WebsoquetService struct {
	Clients map[string]domain.Client
	mu      sync.Mutex
}

// NewWebsoquetService crea una nueva instancia del servicio.
func NewWebsoquetService() *WebsoquetService {
	return &WebsoquetService{
		Clients: make(map[string]domain.Client),
	}
}

// RegisterClient asocia un cliente a un accountID.
func (s *WebsoquetService) RegisterClient(accountID string, client domain.Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Clients[accountID] = client
	log.Printf("Cliente registrado: %s\n", accountID)
}

// SendMessageToAccount envía un mensaje únicamente al cliente asociado a 'receiver'.
func (s *WebsoquetService) SendMessageToAccount(receiver string, msg []byte) {
	s.mu.Lock()
	client, ok := s.Clients[receiver]
	s.mu.Unlock()
	if !ok {
		log.Printf("Cliente %s no encontrado\n", receiver)
		return
	}
	err := client.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Printf("Error enviando mensaje a %s: %v\n", receiver, err)
		s.RemoveClient(receiver)
	}
}

// RemoveClient elimina la conexión de un cliente dado su accountID.
func (s *WebsoquetService) RemoveClient(accountID string) {
	s.mu.Lock()
	client, ok := s.Clients[accountID]
	if ok {
		client.Close()
		delete(s.Clients, accountID)
		log.Printf("Cliente desconectado: %s\n", accountID)
	}
	s.mu.Unlock()
}
