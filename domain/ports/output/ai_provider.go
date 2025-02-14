package outputport

import (
	"context"

	"github.com/josemontano1996/ai-chatbot-backend/domain/entities"
	inputport "github.com/josemontano1996/ai-chatbot-backend/domain/ports/input"
)

// The metadata of the AI service response
type AIResposeMetadata[T any] struct {
	TokensSpent uint32
	Metadata    T
}

type AIProvider[T any] interface {
	SendMessage(ctx context.Context, userMsg *entities.ChatMessage, prevHistory *entities.ChatHistory) (*inputport.AIChatResponse, *AIResposeMetadata[T], error)
	CloseConnection()
}
