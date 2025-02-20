package out

import (
	"time"

	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
)

type TokenAuthService interface {
	GenerateToken(userId string, duration time.Duration) (string, *entities.AuthTokenPayload, error)
	VerifyToken(token string) (*entities.AuthTokenPayload, error)
}
