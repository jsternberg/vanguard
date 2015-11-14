package vanguard

import (
	"encoding/json"
	"net"
	"net/http"
	"sync"

	"code.google.com/p/go-uuid/uuid"

	"github.com/go-martini/martini"
)

type Server struct {
	handler http.Handler
	runs    map[string]*runContext
	runLock sync.RWMutex
}

type runContext struct {
	runner *Runner
	lock   sync.RWMutex
}

func (ctx *runContext) start(plan *Plan) {
	defer ctx.release()
	runner := ctx.runner
	for _, task := range plan.Tasks {
		runner.C <- task
	}
	runner.Close()

	runner.Wait()
}

func (ctx *runContext) release() {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()
	ctx.runner = nil
}

func NewServer() *Server {
	s := &Server{}
	r := martini.NewRouter()
	r.Group("/v1", func(r martini.Router) {
		r.Get("/ping", s.ping)
		r.Post("/provision", s.provision)
		r.Get("/runs/:id/wait", s.wait)
	})

	m := martini.New()
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Action(r.Handle)
	s.handler = m
	s.runs = make(map[string]*runContext)
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

	return http.StatusOK, s.run(plan)
}

func (s *Server) run(plan *Plan) string {
	s.runLock.Lock()
	defer s.runLock.Unlock()

	uuid := uuid.New()
	runner := NewRunner(nil)
	ctx := &runContext{
		runner: runner,
	}
	s.runs[uuid] = ctx

	go ctx.start(plan)
	return uuid
}

func (s *Server) wait(params martini.Params) int {
	uuid := params["id"]

	s.runLock.RLock()
	ctx, ok := s.runs[uuid]
	s.runLock.RUnlock()

	if !ok {
		return http.StatusNotFound
	}
	ctx.runner.Wait()
	return http.StatusOK
}
