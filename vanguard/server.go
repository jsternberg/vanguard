package vanguard

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/go-martini/martini"
)

type Server struct {
	handler http.Handler
}

func NewServer() *Server {
	s := &Server{}
	r := martini.NewRouter()
	r.Get("/v1/ping", s.Ping)

	m := martini.New()
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Action(r.Handle)
	s.handler = m
	return s
}

func (s *Server) Serve(l net.Listener) error {
	return http.Serve(l, s.handler)
}

func (s *Server) Ping() []byte {
	serverInfo := map[string]string{
		"version": Version,
	}
	out, _ := json.MarshalIndent(&serverInfo, "", "  ")
	return out
}
