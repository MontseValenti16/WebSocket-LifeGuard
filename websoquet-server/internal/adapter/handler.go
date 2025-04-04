package adapter

import (
	"encoding/json"
	"log"
	"net/http"
	"WS/websoquet-server/internal/domain"
	"WS/websoquet-server/internal/service"

	"github.com/gorilla/websocket"
)

// Handler gestiona las conexiones WebSocket.
type Handler struct {
	Service *service.WebsoquetService
}

// NewHandler crea una nueva instancia de Handler.
func NewHandler(svc *service.WebsoquetService) *Handler {
	return &Handler{Service: svc}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ServeWS actualiza la conexión HTTP a WebSocket, registra al cliente y lo asocia a una cuenta.
func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al actualizar a WebSocket:", err)
		return
	}
	client := NewClient(conn)

	// Se espera que la URL incluya ?account=ID
	accountID := r.URL.Query().Get("account")
	if accountID == "" {
		log.Println("No se proporcionó el accountID en la URL")
		client.Close()
		return
	}

	h.Service.RegisterClient(accountID, client)
	go h.handleMessages(accountID, client)
}

// handleMessages procesa los mensajes recibidos del cliente.
// handleMessages procesa los mensajes recibidos del cliente.
func (h *Handler) handleMessages(accountID string, client *Client) {
	defer func() {
		h.Service.RemoveClient(accountID, client) // 🔥 Se pasa el cliente correcto
	}()
	for {
		_, msg, err := client.ReadMessage()
		if err != nil {
			log.Println("Error leyendo mensaje:", err)
			break
		}
		var message domain.Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Error decodificando JSON:", err)
			continue
		}

		// Aquí mostramos tanto los campos fijos como los adicionales.
		log.Printf("Mensaje recibido de %s para %s con datos adicionales: %+v\n", message.Sender, message.Receiver, message.Data)

		// Extraer valores individuales de message.Data
		if dataMap, ok := message.Data["content"].(map[string]interface{}); ok {
			accion, _ := dataMap["accion"].(string)
			fecha, _ := dataMap["fecha"].(string)
			idPersona, _ := dataMap["id_persona"].(string)
			macAddress, _ := dataMap["macaddress"].(string)
			mensaje, _ := dataMap["mensaje"].(string)
			nivelPeligro, _ := dataMap["nivel_peligro"].(string)

			log.Println("Acción:", accion)
			log.Println("Fecha:", fecha)
			log.Println("ID Persona:", idPersona)
			log.Println("MAC Address:", macAddress)
			log.Println("Mensaje:", mensaje)
			log.Println("Nivel de Peligro:", nivelPeligro)
		} else {
			log.Println("No se encontró 'content' en los datos adicionales.")
		}

		msgToSend, err := json.Marshal(message)
		if err != nil {
			log.Println("Error codificando JSON:", err)
			continue
		}
		h.Service.SendMessageToAccount(message.Receiver, msgToSend)
	}
}



	

