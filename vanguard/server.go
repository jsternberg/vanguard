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
	r.Get("/v1/ping", s.ping)
	r.Post("/v1/provision", s.provision)

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

func (s *Server) ping() []byte {
	serverInfo := map[string]string{
		"version": Version,
	}
	out, _ := json.MarshalIndent(&serverInfo, "", "  ")
	return out
}

func (s *Server) provision(req *http.Request) (int, string) {
	defer req.Body.Close()

	loader := YamlLoader{}
	plan, err := loader.Load(req.Body)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	runner := NewRunner(nil)
	for _, task := range plan.Tasks {
		runner.C <- task
	}
	runner.Close()

	runner.Wait()
	return http.StatusOK, ""
}
