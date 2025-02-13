package services

import (
	"github.com/gin-gonic/gin"
	"github.com/josemontano1996/ai-chatbot-backend/sharedtypes"
)



type ChatResponse struct {
	AIResponse     *sharedtypes.Message
	UpdatedHistory *sharedtypes.History
	TokenCount     int32
}

type AIService interface {
	Chat(ctx *gin.Context, userMessage *sharedtypes.Message, prevHistory *sharedtypes.History) (*ChatResponse, error)
}
