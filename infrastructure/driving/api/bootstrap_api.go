package api

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	geminiadapter "github.com/josemontano1996/ai-chatbot-backend/infrastructure/driven/ai_providers/gemini"
	"github.com/josemontano1996/ai-chatbot-backend/infrastructure/driven/auth"
	sqlcrepo "github.com/josemontano1996/ai-chatbot-backend/infrastructure/driven/repository/db"
	redisrepo "github.com/josemontano1996/ai-chatbot-backend/infrastructure/driven/repository/redis"

	controller "github.com/josemontano1996/ai-chatbot-backend/infrastructure/driving/api/controllers"
	"github.com/josemontano1996/ai-chatbot-backend/infrastructure/driving/ws"
	chatws "github.com/josemontano1996/ai-chatbot-backend/infrastructure/driving/ws/chat"
	"github.com/josemontano1996/ai-chatbot-backend/internal/config"
	"github.com/josemontano1996/ai-chatbot-backend/internal/usecases"
)

func RunRestApi() {
	config, err := config.LoadEnv("./", "prod")

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Infraestructure layer setup
	// WS
	wsConfig, err := ws.NewWSConfig(1024, 1024, 30*time.Minute)
	if err != nil {
		log.Fatalf("failed to create WS config: %v", err)
	}
	AIWSChatClient, err := chatws.NewAIChatWSClient(*wsConfig)
	if err != nil {
		log.Fatalf("failed to create WS instance: %v", err)
	}

	// Repository layer setup
	// DB store
	conn, err := pgxpool.New(context.Background(), config.PostgresConnectionString)
	if err != nil {
		log.Fatalf("failed to create connection pool: %v", err)
	}
	defer conn.Close()
	userRepository := sqlcrepo.NewUserRepository(conn)

	// KV store
	redisConfig, err := redisrepo.NewRedisConfig(config.RedisAddress, config.RedisPassword, config.RedisDB, 30*time.Minute)
	if err != nil {
		log.Fatalf("failed to create Redis config: %v", err)
	}
	redisRepo, err := redisrepo.NewRedisRepository(redisConfig)
	if err != nil {
		log.Fatalf("failed to create Redis repository: %v", err)
	}
	defer redisRepo.Close()

	// AIProviders
	geminiConfig, err := geminiadapter.NewGeminiConfig(config.GeminiApiKey, geminiadapter.Gemini15FlashModelName, int32(config.GeminiMaxTokens))

	if err != nil {
		log.Fatalf("failed to create Gemini config: %v", err)
	}

	geminiProvider, err := geminiadapter.NewGeminiAdapter(context.Background(), *geminiConfig)

	if err != nil {
		log.Fatalf("failed to create Gemini service: %v", err)
	}
	defer geminiProvider.CloseConnection()

	// AuthService
	pasetoAuthService := auth.NewPasetoAuthenticator()

	// Domain layer setup
	AIChatUseCase := usecases.NewAIChatUseCase(geminiProvider, redisRepo)
	AuthUseCase, err := usecases.NewAuthUseCase(pasetoAuthService, userRepository, config.SessionDuration)

	if err != nil {
		log.Fatalf("error in auth use case: %v", err)
	}

	// Interface/Presenter layer setup
	AIController := controller.NewAIController(AIChatUseCase, redisRepo, AIWSChatClient)
	AuthController := controller.NewAuthController(AuthUseCase, userRepository)

	// Create Gin Router and register routes
	server := NewServer(&config)
	server.RegisterRoutes(&AuthUseCase, AuthController, AIController)

	// Start server
	err = server.RunServer(config.ServerPort)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
