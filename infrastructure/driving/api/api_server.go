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
	"github.com/josemontano1996/ai-chatbot-backend/internal/config"
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/in"
	"github.com/rs/cors"
)

type Server struct {
	router *gin.Engine
	srv    *http.Server
	env    *config.Env
}

func NewApiServer(env *config.Env) *Server {
	return &Server{
		router: gin.Default(),
		env:    env,
	}
}

func (s *Server) RegisterRoutes(authUseCases *in.AuthUseCase, authController *controller.AuthController, AIController *controller.AIController, userController *controller.UserController) {
	apiRoutes := s.router.Group("/api")
	{
		apiRoutes.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"status": "healthy"})
		})

		apiRoutes.POST("/register", authController.RegisterUser)
		apiRoutes.POST("/login", authController.Login)

		privateGroup := apiRoutes.Group("/private")

		privateGroup.Use(middleware.AuthMiddleware(*authUseCases))
		{
			userRoutes := privateGroup.Group("/user")
			userRoutes.GET("/", userController.GetUserById)

			wsRoutes := privateGroup.Group("/ws")
			wsRoutes.GET("/chat", AIController.ChatWithAI)
		}

	}

}

func (s *Server) RunServer(port string) error {
	isProd := true
	if s.env.AppEnvironment == "dev" {
		isProd = false
	}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{s.env.FrontEndOrigin, "http://localhost:4173", " http://localhost:5173"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
			"Upgrade",                  //  MUST include "Upgrade" for WebSocket
			"Connection",               //  MUST include "Connection" for WebSocket
			"Sec-WebSocket-Key",        //  MUST include "Sec-WebSocket-Key" for WebSocket
			"Sec-WebSocket-Version",    //  MUST include "Sec-WebSocket-Version" for WebSocket
			"Origin",                   //  Good practice to include, often needed
			"Sec-WebSocket-Extensions", // Include if you are using WebSocket extensions (like permessage-deflate)
		},
		AllowCredentials: true,
		Debug:            isProd,
	})

	handler := c.Handler(s.router)

	s.srv = &http.Server{ // Initialize the HTTP server
		Addr:    ":" + port,
		Handler: handler,
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
