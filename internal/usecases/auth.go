package usecases

import (
	"context"
	"fmt"

	"github.com/josemontano1996/ai-chatbot-backend/internal/dto"
	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/in"
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/out"
)

type AuthUseCases struct {
	pasetoService  out.PasetoService
	userRepository out.UserRepository
}

func NewAuthUseCases(ps out.PasetoService, ur out.UserRepository) in.AuthService {
	return &AuthUseCases{
		pasetoService:  ps,
		userRepository: ur,
	}
}

func (uc *AuthUseCases) Login(ctx context.Context, email, password string) (string, error) {
	user, err := uc.userRepository.FindByEmail(ctx, email)

	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}

	// TODO: implement comparing password hashing
	if user == nil || user.Password != password {
		return "", fmt.Errorf("login failed: invalid credentials")
	}

	token, err := uc.pasetoService.GenerateToken(user)

	if err != nil {
		return "", fmt.Errorf("login failed: falided to generate token: %w", err)
	}

	return token, nil
}

func (uc *AuthUseCases) ValidateToken(ctx context.Context, token string) (*dto.User, error) {
	user, err := uc.pasetoService.ValidateToken(token)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return dto.NewUserDTOFromEntity(user)
}

func (uc *AuthUseCases) GenerateToken(user *entities.User) (string, error) {
	return uc.pasetoService.GenerateToken(user)
}
