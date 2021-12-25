package internalhttp

import (
	"context"
	"net"
	"net/http"
	"strconv"
)

type Server struct {
	host   string
	port   int
	logger Logger
	server *http.Server
}

type Logger interface {
	Info(msg string, params ...interface{})
	LogHTTPRequest(r *http.Request, code, length int)
}

type Application interface{}

func NewServer(logger Logger, app Application, host string, port int) *Server {
	myServer := &Server{
		host:   host,
		port:   port,
		logger: logger,
		server: nil,
	}

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(host, strconv.Itoa(port)),
		Handler: loggingMiddleware(http.HandlerFunc(myServer.handleHTTP), logger),
	}

	myServer.server = httpServer

	return myServer
}

func (s *Server) handleHTTP(w http.ResponseWriter, r *http.Request) {
	msg := []byte("Hello, world!\n")
	w.WriteHeader(200)
	w.Write(msg)
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("Start HTTP Server", "host", s.host, "port", s.port)
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
