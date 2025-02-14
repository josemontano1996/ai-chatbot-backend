package inputport

import (
	"github.com/josemontano1996/ai-chatbot-backend/hex/domain/entities"
)

// ChatResponse encapsulates the domain's response to a chat request
type AIChatResponse struct {
	ChatMessage *entities.ChatMessage
}

// ChatService defines the use cases related to chat functionality
type AIChatService interface {
	SendChatMessage(userMessage *entities.ChatMessage, history *entities.ChatHistory) (*AIChatResponse, error)
}
