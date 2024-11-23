package webserver

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type WebServer struct {
	WebServerPort string
	Router        *chi.Mux
	Server        *http.Server
	Handlers      map[string]struct {
		Handler http.HandlerFunc
		Method  string
	}
}

func NewWebServer(port string) *WebServer {
	return &WebServer{
		WebServerPort: port,
		Router:        chi.NewRouter(),
		Handlers: make(map[string]struct {
			Handler http.HandlerFunc
			Method  string
		}),
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc, method string) {
	key := path + "_" + method
	s.Handlers[key] = struct {
		Handler http.HandlerFunc
		Method  string
	}{Handler: handler, Method: method}
}

func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)

	s.Server = &http.Server{
		Addr:    s.WebServerPort,
		Handler: s.Router,
	}

	for key, entry := range s.Handlers {
		switch entry.Method {
		case "GET":
			s.Router.Get(key[:len(key)-len("_GET")], entry.Handler)
		case "POST":
			s.Router.Post(key[:len(key)-len("_POST")], entry.Handler)
		case "PUT":
			s.Router.Put(key[:len(key)-len("_PUT")], entry.Handler)
		case "PATCH":
			s.Router.Patch(key[:len(key)-len("_PATCH")], entry.Handler)
		case "DELETE":
			s.Router.Delete(key[:len(key)-len("_DELETE")], entry.Handler)
		default:
			s.Router.Method(entry.Method, key[:len(key)-len("_"+entry.Method)], entry.Handler)
		}
	}

	go func() {
		if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}

func (s *WebServer) Stop() error {
	if s.Server != nil {
		return s.Server.Close()
	}
	return nil
}
