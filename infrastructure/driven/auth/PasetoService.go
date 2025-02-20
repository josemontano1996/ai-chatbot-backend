package auth

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/out"
)

type PasetoAuthenticator struct {
	simmetricKey paseto.V4SymmetricKey
}

func NewPasetoAuthenticator() out.TokenAuthService {
	return &PasetoAuthenticator{
		simmetricKey: paseto.NewV4SymmetricKey(),
	}
}

func (p *PasetoAuthenticator) GenerateToken(userId string, duration time.Duration) (string, *entities.AuthTokenPayload, error) {
	// create paseto token
	token := paseto.NewToken()

	payload, err := entities.NewAuthTokenPayload(userId, duration)

	if err != nil {
		return "", payload, err
	}

	token.Set("id", payload.ID)
	token.Set("user_id", payload.UserId)
	token.SetIssuedAt(payload.IssuedAt)
	token.SetExpiration(payload.ExpiredAt)

	return token.V4Encrypt(p.simmetricKey, nil), payload, nil
}

func (p *PasetoAuthenticator) VerifyToken(token string) (*entities.AuthTokenPayload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())

	parsedToken, err := parser.ParseV4Local(p.simmetricKey, token, nil)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return getPayloadFromToken(parsedToken)
}

func getPayloadFromToken(t *paseto.Token) (*entities.AuthTokenPayload, error) {

	id, err := t.GetString("id")
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	userId, err := t.GetString("user_id")
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	issuedAt, err := t.GetIssuedAt()
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	expiredAt, err := t.GetExpiration()
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return entities.CreateAuthTokenPayload(id, userId, issuedAt, expiredAt)
}
