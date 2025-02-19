package in

import (
	"context"

	"github.com/josemontano1996/ai-chatbot-backend/internal/dto"
	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	ValidateToken(ctx context.Context, token string) (*dto.User, error)
	GenerateToken(user *entities.User) (string, error)
}
