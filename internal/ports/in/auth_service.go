package in

import (
	"context"
	"time"

	"github.com/josemontano1996/ai-chatbot-backend/internal/dto"
	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
)

type AuthUseCase interface {
	RegisterUser(ctx context.Context, email, password string) (*dto.User, error)
	Login(ctx context.Context, email, password string) (string, *entities.AuthTokenPayload, error)
	ValidateToken(ctx context.Context, token string) (*entities.AuthTokenPayload, error)
	GenerateToken(userId string, tokenDuration time.Duration) (string, *entities.AuthTokenPayload, error)
}
