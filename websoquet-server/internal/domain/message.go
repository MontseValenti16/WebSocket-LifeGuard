package domain

// Message representa la estructura de un mensaje.
type Message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Content  string `json:"content"`
}

// Client define las operaciones básicas que debe implementar un cliente WebSocket.
type Client interface {
	ReadMessage() (int, []byte, error)
	WriteMessage(messageType int, msg []byte) error
	Close() error
}
