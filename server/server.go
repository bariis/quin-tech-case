package server

import (
	"net/http"
	"time"

	"github.com/bariis/quin-tech-case/task"
	"github.com/gorilla/mux"
)

// Server represents an HTTP server. It wraps all HTTP functionality used by the application.
type Server struct {
	Server *http.Server

	TaskService task.TaskService
}

// NewServer returns a new instance of Server.
func NewServer() *Server {
	r := mux.NewRouter()
	s := &Server{
		Server: &http.Server{
			Addr:         ":5000",
			Handler:      r,
			WriteTimeout: 5 * time.Second,
			ReadTimeout:  5 * time.Second,
		},
	}

	s.registerTaskRoutes(r)
	return s
}
