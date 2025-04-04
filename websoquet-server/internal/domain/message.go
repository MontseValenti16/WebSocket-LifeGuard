package domain

import "encoding/json"

// Message representa un mensaje con campos fijos y datos adicionales arbitrarios.
type Message struct {
	Sender   string                 `json:"sender"`
	Receiver string                 `json:"receiver"`
	// Content se mantiene con el valor recibido (puede ser string u objeto).
	Content interface{}            `json:"content,omitempty"`
	// Data almacena los demás campos que se envíen en el JSON.
	Data    map[string]interface{} `json:"-"`
}

// UnmarshalJSON implementa un decodificador personalizado para extraer sender, receiver y content.
func (m *Message) UnmarshalJSON(data []byte) error {
	// Decodificar en un mapa genérico.
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if s, ok := raw["sender"].(string); ok {
		m.Sender = s
	}
	if r, ok := raw["receiver"].(string); ok {
		m.Receiver = r
	}
	// Si se envía "content", se asigna tal cual.
	if c, ok := raw["content"]; ok {
		m.Content = c
		delete(raw, "content")
	}

	// Eliminar los campos fijos para que en Data queden solo los adicionales.
	delete(raw, "sender")
	delete(raw, "receiver")
	m.Data = raw

	return nil
}

// La interfaz Client se mantiene igual.
type Client interface {
	ReadMessage() (int, []byte, error)
	WriteMessage(messageType int, msg []byte) error
	Close() error
}
