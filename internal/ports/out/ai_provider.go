package out

import (
	"context"

	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/in"
)

// The metadata of the AI service response
type AIResposeMetadata[T any] struct {
	TokensSpent uint32
	Metadata    T
}

type AIProvider[T any] interface {
	SendMessage(ctx context.Context, userMsg *entities.ChatMessage, prevHistory *entities.ChatHistory) (*in.AIChatResponse, *AIResposeMetadata[T], error)
	CloseConnection()
}
