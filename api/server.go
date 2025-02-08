package api

import "github.com/gin-gonic/gin"

type Server struct {
	router *gin.Engine
}

func NewServer() (*Server, error) {

	server := &Server{
		router: gin.Default(),
	}
	
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
