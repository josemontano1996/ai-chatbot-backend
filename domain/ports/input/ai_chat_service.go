package inputport

import (
	"context"

	"github.com/josemontano1996/ai-chatbot-backend/domain/entities"
)

// ChatResponse encapsulates the domain's response to a chat request
type AIChatResponse struct {
	ChatMessage *entities.ChatMessage
}

// ChatService defines the use cases related to chat functionality
type AIChatUseCase interface {
	SendChatMessage(ctx context.Context, user *entities.User, userMessage *entities.ChatMessage, history *entities.ChatHistory) (*AIChatResponse, error)
}
