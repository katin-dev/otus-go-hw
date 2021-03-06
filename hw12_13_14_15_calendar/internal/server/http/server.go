package internalhttp

import (
	"context"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/katin.dev/otus-go-hw/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	host   string
	port   int
	logger Logger
	server *http.Server
}

type Logger interface {
	Info(msg string, params ...interface{})
	Error(msg string, params ...interface{})
	LogHTTPRequest(r *http.Request, code, length int)
}

func NewServer(logger Logger, app *app.App, host string, port int) *Server {
	myServer := &Server{
		host:   host,
		port:   port,
		logger: logger,
		server: nil,
	}

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(host, strconv.Itoa(port)),
		Handler: loggingMiddleware(NewRouter(app, logger), logger),
	}

	myServer.server = httpServer

	return myServer
}

func NewRouter(app *app.App, log Logger) http.Handler {
	handlers := NewServerHandlers(app, log)

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HelloWorld).Methods("GET")
	r.HandleFunc("/events", handlers.CreateEvent).Methods("POST")
	r.HandleFunc("/events/{id}", handlers.UpdateEvent).Methods("PUT")
	r.HandleFunc("/events/{id}", handlers.DeleteEvent).Methods("DELETE")
	r.HandleFunc("/events", handlers.ListEvents).Methods("GET")

	return r
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("Start HTTP Server on %s:%d", s.host, s.port)
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.server.Shutdown(ctx)
	return nil
}
