package in

import (
	"context"

	"github.com/josemontano1996/ai-chatbot-backend/internal/dto"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	ValidateToken(ctx context.Context, token string) (*dto.User, error)
	GenerateToken(userId string) (string, error)
}
