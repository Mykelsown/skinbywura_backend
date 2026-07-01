package server

import (
	"net/http"

	"github.com/Mykelsown/skinbywura_backend.git/config"
)

type Server struct {
	address string
	route   http.Handler
}

func New(cfg *config.EnvData) *Server {	
	ser := &Server {
		address: cfg.Port,
		route: loadRoute(),
	}

	return ser
}

func (s *Server) Run() error {
	// fmt.Println("binding to: ", s.address)
	err :=http.ListenAndServe(":"+s.address, s.route)
	return err
}