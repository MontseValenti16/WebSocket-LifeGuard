package domain

import "encoding/json"

// Message representa un mensaje con campos fijos y datos adicionales arbitrarios.
type Message struct {
    Sender   string                 `json:"sender"`
    Receiver string                 `json:"receiver"`
    Content  string                 `json:"content,omitempty"` // Nuevo campo
    Data     map[string]interface{} `json:"-"`
}

// UnmarshalJSON implementa un decodificador personalizado para extraer sender, receiver y content.
func (m *Message) UnmarshalJSON(data []byte) error {
    // Decodificar en un mapa
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
    if c, ok := raw["content"].(string); ok {
        m.Content = c
        delete(raw, "content")
    }

    delete(raw, "sender")
    delete(raw, "receiver")
    m.Data = raw

    return nil
}

// Client interface permanece igual
type Client interface {
    ReadMessage() (int, []byte, error)
    WriteMessage(messageType int, msg []byte) error
    Close() error
}
