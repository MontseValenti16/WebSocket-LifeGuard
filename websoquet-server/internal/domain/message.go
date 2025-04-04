package domain

import "encoding/json"

// Message representa un mensaje con campos fijos y datos adicionales arbitrarios.
type Message struct {
	Sender   string                 `json:"sender"`
	Receiver string                 `json:"receiver"`
	// Content es un campo opcional; si el JSON lo envía como string, se asigna aquí.
	Content string `json:"content,omitempty"`
	// Data contendrá todos los campos adicionales que no se pudieron mapear a los campos fijos.
	Data map[string]interface{} `json:"-"`
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
	// Si "content" es un string, lo asignamos a Content y lo eliminamos del mapa.
	if c, ok := raw["content"].(string); ok {
		m.Content = c
		delete(raw, "content")
	}
	// Los demás campos se guardan en Data.
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
