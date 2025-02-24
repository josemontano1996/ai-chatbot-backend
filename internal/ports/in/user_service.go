package in

import (
	"context"

	"github.com/google/uuid"
	"github.com/josemontano1996/ai-chatbot-backend/internal/dto"
)

type UserUseCase interface {
	GetUserById(ctx context.Context, id uuid.UUID) (*dto.User, error)
}
