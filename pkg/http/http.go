package http

import (
	"moe.two.bgmi-gls/pkg/http/m3u8"
	"moe.two.bgmi-gls/pkg/http/middleware"
	"moe.two.bgmi-gls/pkg/http/ts"
	"net/http"
)
import "github.com/gorilla/mux"
import log "github.com/sirupsen/logrus"

type Server struct {
	port string
}

func NewServer(port string) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) ListenAndServe() error {
	registryRouter()
	log.Infof("Starting Http Server On %s.", s.port)
	return http.ListenAndServe(s.port, nil)
}

func registryRouter(){
	router := mux.NewRouter()
	router.HandleFunc("/ts/{path:.*$}", ts.VideoFragment)
	router.HandleFunc("/m3u8/{path:.*$}", m3u8.VideoIndex)
	router.Use(middleware.LoggingMiddleware)

	http.Handle("/", router)
}
