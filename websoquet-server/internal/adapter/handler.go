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
func (h *Handler) handleMessages(accountID string, client *Client) {
	defer func() {
		h.Service.RemoveClient(accountID, client)
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

		log.Printf("Mensaje recibido de %s para %s\n", message.Sender, message.Receiver)
		// Primero se revisa si Content fue asignado (en caso de enviarlo como string)
		if message.Content != "" {
			log.Println("Content:", message.Content)
		} else {
			// Si no se asignó Content, se intenta extraer "content" desde Data.
			if dataMap, ok := message.Data["content"].(map[string]interface{}); ok {
				accion, _ := dataMap["accion"].(string)
				fecha, _ := dataMap["fecha"].(string)
				// id_persona puede ser numérico o string, se captura como interface{}
				idPersona := dataMap["id_persona"]
				macAddress, _ := dataMap["macaddress"].(string)
				mensajeText, _ := dataMap["mensaje"].(string)
				nivelPeligro, _ := dataMap["nivel_peligro"].(string)

				log.Println("Acción:", accion)
				log.Println("Fecha:", fecha)
				log.Println("ID Persona:", idPersona)
				log.Println("MAC Address:", macAddress)
				log.Println("Mensaje:", mensajeText)
				log.Println("Nivel de Peligro:", nivelPeligro)
			} else {
				log.Println("No se encontró 'content' en el mensaje.")
			}
		}

		// Vuelve a serializar el mensaje para enviarlo.
		msgToSend, err := json.Marshal(message)
		if err != nil {
			log.Println("Error codificando JSON:", err)
			continue
		}
		h.Service.SendMessageToAccount(message.Receiver, msgToSend)
	}
}
