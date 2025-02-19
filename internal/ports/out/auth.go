package out

import (
	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
)

type PasetoService interface {
	GenerateToken(userId string) (string, error)
	ValidateToken(token string) (*entities.User, error)
}
