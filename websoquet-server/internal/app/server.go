package app

import (
	"net/http"
	"WS/websoquet-server/internal/adapter"
	"WS/websoquet-server/internal/service"
)

type Server struct {
	Handler *adapter.Handler
}

func NewServer() *Server {
	svc := service.NewWebsoquetService()
	handler := adapter.NewHandler(svc)
	return &Server{
		Handler: handler,
	}
}

func (s *Server) Start(addr string) error {
	http.HandleFunc("/ws", s.Handler.ServeWS)
	return http.ListenAndServe(addr, nil)
}
