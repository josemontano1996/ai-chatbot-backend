package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/josemontano1996/ai-chatbot-backend/infrastructure/driving/api"
)

type Server struct {
	router *gin.Engine
	srv    *http.Server
}

func NewServer() *Server {
	return &Server{
		router: gin.Default(),
	}
}

func (s *Server) RegisterRoutes(AIController *api.AIController) {

	s.router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	s.router.GET("/chat", AIController.ChatWithAI)
}

func (s *Server) RunServer(port string) error {
	s.srv = &http.Server{ // Initialize the HTTP server
		Addr:    ":" + port,
		Handler: s.router,
	}

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
	return nil // Return nil on successful shutdown
}
