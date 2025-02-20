package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	controller "github.com/josemontano1996/ai-chatbot-backend/infrastructure/driving/api/controllers"
	"github.com/josemontano1996/ai-chatbot-backend/infrastructure/driving/api/middleware"
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/in"
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

func (s *Server) RegisterRoutes(AuthUseCases *in.AuthUseCase, AuthController *controller.AuthController, AIController *controller.AIController) {
	apiRoutes := s.router.Group("/api")
	{
		apiRoutes.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"status": "healthy"})
		})

		privateGroup := apiRoutes.Group("/private")
		privateGroup.Use(middleware.AuthMiddleware(*AuthUseCases))
		{
		}

		s.router.GET("/chat", AIController.ChatWithAI)
	}

}

func (s *Server) RunServer(port string) error {
	s.srv = &http.Server{ // Initialize the HTTP server
		Addr:    ":" + port,
		Handler: s.router,
	}
	fmt.Println("Server running on port: ", port)
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
