package outrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email, password string) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (userEntity *entities.User, password string, err error)
	FindById(ctx context.Context, id uuid.UUID) (userEntity *entities.User, password string, err error)
	UpdateUser(ctx context.Context, id uuid.UUID, new_email, new_password string) (*entities.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
