package api

import (
	"github.com/gin-gonic/gin"
	"github.com/josemontano1996/ai-chatbot-backend/api/routes"
)

type Server struct {
	router *gin.Engine
}

func NewServer() (*Server, error) {
	router, err := routes.RegisterRoutes(gin.Default())

	if err != nil {
		return &Server{}, err
	}

	server := &Server{
		router: router,
	}

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
