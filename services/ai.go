package services

import (
	"github.com/gin-gonic/gin"
	"github.com/josemontano1996/ai-chatbot-backend/sharedtypes"
)

const (
	AIUserMessageCode   int8 = 1
	AIBotResponseCode   int8 = 2
	AISystemMessageCode int8 = 3
)

type ChatResponse struct {
	AIResponse     *sharedtypes.Message
	// UpdatedHistory *sharedtypes.History
	TokenCount     int32
}


type AIService interface {
	Chat(ctx *gin.Context, userMessage *sharedtypes.Message, prevHistory *sharedtypes.History) (*ChatResponse, error)
}
