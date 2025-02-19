package out

import (
	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
)

type PasetoService interface {
	GenerateToken(user *entities.User) (string, error)
	ValidateToken(token string) (*entities.User, error)
}
