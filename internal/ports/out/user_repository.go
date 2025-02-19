package out

import (
	"context"

	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
}
