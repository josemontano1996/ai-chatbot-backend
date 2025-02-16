package out

import (
	"context"

	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
)

type ChatMessageRepository interface {
	GetChatHistory(ctx context.Context, key string) (*entities.ChatHistory, error)
	SaveMessage(ctx context.Context, key string, message *entities.ChatMessage) error
	SaveMessages(ctx context.Context, key string, messages ...*entities.ChatMessage) error
}
