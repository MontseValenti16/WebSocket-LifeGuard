package service

import (
	"log"
	"sync"
	"WS/websoquet-server/internal/domain"

	"github.com/gorilla/websocket"
)

type WebsoquetService struct {
	Clients map[string][]domain.Client // Ahora es una lista de clientes por mac address
	mu      sync.Mutex
}

func NewWebsoquetService() *WebsoquetService {
	return &WebsoquetService{
		Clients: make(map[string][]domain.Client),
	}
}


// RegisterClient asocia un cliente a un accountID.
func (s *WebsoquetService) RegisterClient(accountID string, client domain.Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Agrega el nuevo cliente a la lista existente de la misma mac address
	s.Clients[accountID] = append(s.Clients[accountID], client)
	log.Printf("Cliente registrado: %s\n", accountID)
}


// SendMessageToAccount envía un mensaje únicamente al cliente asociado a 'receiver'.
func (s *WebsoquetService) SendMessageToAccount(receiver string, msg []byte) {
	s.mu.Lock()
	clients, ok := s.Clients[receiver]
	s.mu.Unlock()
	if !ok {
		log.Printf("Cliente %s no encontrado\n", receiver)
		return
	}

	// Enviar mensaje a todos los clientes conectados con la misma mac address
	for _, client := range clients {
		err := client.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("Error enviando mensaje a %s: %v\n", receiver, err)
			s.RemoveClient(receiver, client) // Remover solo el cliente con error
		}
	}
}


// RemoveClient elimina la conexión de un cliente dado su accountID.
func (s *WebsoquetService) RemoveClient(accountID string, clientToRemove domain.Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	clients := s.Clients[accountID]
	newClients := []domain.Client{}

	// Filtramos la lista eliminando solo el cliente que se desconectó
	for _, client := range clients {
		if client != clientToRemove {
			newClients = append(newClients, client)
		} else {
			client.Close()
			log.Printf("Cliente desconectado de %s\n", accountID)
		}
	}

	// Si ya no hay clientes conectados con esa mac address, eliminamos la clave
	if len(newClients) == 0 {
		delete(s.Clients, accountID)
	} else {
		s.Clients[accountID] = newClients
	}
}

