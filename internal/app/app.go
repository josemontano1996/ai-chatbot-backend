package app

import (
	"context"
	"log"

	"github.com/josemontano1996/ai-chatbot-backend/config"
	"github.com/josemontano1996/ai-chatbot-backend/domain/usecases"
	geminiadapter "github.com/josemontano1996/ai-chatbot-backend/infrastructure/driven/ai_providers/gemini"
	repository "github.com/josemontano1996/ai-chatbot-backend/infrastructure/driven/respository/redis"
	api "github.com/josemontano1996/ai-chatbot-backend/infrastructure/driving/api/controllers"
	"github.com/josemontano1996/ai-chatbot-backend/infrastructure/server"
)

func StartApp() {
	config, err := config.LoadEnv("./", "prod")

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Infraestructure layer setup
	redisConfig := repository.NewRedisConfig(config.RedisAddress, config.RedisPassword, config.RedisDB)
	redisRepo := repository.NewRedisRepository(redisConfig)
	defer redisRepo.Close()

	geminiConfig, err := geminiadapter.NewGeminiConfig(config.GeminiApiKey, geminiadapter.Gemini15FlashModelName, int32(config.GeminiMaxTokens))

	if err != nil {
		log.Fatalf("failed to create Gemini config: %v", err)
	}

	geminiProvider, err := geminiadapter.NewGeminiAdapter(context.Background(), *geminiConfig)

	if err != nil {
		log.Fatalf("failed to create Gemini service: %v", err)
	}
	defer geminiProvider.CloseConnection()

	// Domain layer setup
	AIChatUseCase := usecases.NewAIChatUseCase(geminiProvider, redisRepo)

	// Interface/Presenter layer setup
	AIController := api.NewAIController(AIChatUseCase, redisRepo)

	// Create Gin Router and register routes
	server := server.NewServer()
	server.RegisterRoutes(AIController)

	// Start server
	err = server.RunServer(config.ServerPort)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
