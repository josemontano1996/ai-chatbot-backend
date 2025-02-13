package services

import (
	"github.com/josemontano1996/ai-chatbot-backend/sharedtypes"
)

type ChatResponse struct {
	AIResponse       *sharedtypes.Message
	TotalTokensSpend uint32
}

type AIService[T any] interface {
	SendChatMessage(userMessage *sharedtypes.Message, prevHistory *sharedtypes.History) (response *ChatResponse, metadata T, err error)
	Close()
}
