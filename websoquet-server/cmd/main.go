package main

import (
	"fmt"
	"log"
	"WS/websoquet-server/internal/app"
)

func main() {
	server := app.NewServer()
	addr := ":8080"
	fmt.Printf("Servidor WebSocket ejecutándose en %s\n", addr)
	if err := server.Start(addr); err != nil {
		log.Fatal("Error en el servidor:", err)
	}
}
