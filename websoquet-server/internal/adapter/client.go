package adapter

import (
	"log"
	"WS/websoquet-server/internal/domain"

	"github.com/gorilla/websocket"
)

// Client representa un cliente conectado vía WebSocket.
type Client struct {
	Conn *websocket.Conn
}

// NewClient crea una nueva instancia de Client.
func NewClient(conn *websocket.Conn) *Client {
	return &Client{Conn: conn}
}

// ReadMessage lee un mensaje del cliente.
func (c *Client) ReadMessage() (int, []byte, error) {
	return c.Conn.ReadMessage()
}

// WriteMessage envía un mensaje al cliente.
func (c *Client) WriteMessage(messageType int, msg []byte) error {
	err := c.Conn.WriteMessage(messageType, msg)
	if err != nil {
		log.Println("Error enviando mensaje:", err)
	}
	return err
}

// Close cierra la conexión del cliente.
func (c *Client) Close() error {
	return c.Conn.Close()
}

// Garantizamos que Client implementa la interfaz domain.Client.
var _ domain.Client = (*Client)(nil)
