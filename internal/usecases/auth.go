package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/josemontano1996/ai-chatbot-backend/internal/dto"
	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/in"
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/out"
	outrepo "github.com/josemontano1996/ai-chatbot-backend/internal/ports/out/repositories"
	"github.com/josemontano1996/ai-chatbot-backend/pkg/utils"
)

type AuthUseCases struct {
	auth           out.TokenAuthService
	userRepository outrepo.UserRepository
	tokenDuration  time.Duration
}

func NewAuthUseCase(ps out.TokenAuthService, ur outrepo.UserRepository, tokenDuration time.Duration) (in.AuthUseCase, error) {

	if tokenDuration == 0 {
		return nil, errors.New("token duration must be greater than 0")
	}

	return &AuthUseCases{
		auth:           ps,
		userRepository: ur,
		tokenDuration:  tokenDuration,
	}, nil
}

func (uc *AuthUseCases) RegisterUser(ctx context.Context, email, password string) (*dto.User, error) {

	if len(password) < 8 {
		return nil, fmt.Errorf("invalid user password, length is lower than 8 chars")
	}

	hashedPassword, err := utils.HashPassword(password, 12)
	if err != nil {
		return nil, err
	}

	userEntity, err := uc.userRepository.CreateUser(ctx, email, hashedPassword)
	if err != nil {
		return nil, err
	}

	return dto.NewUserDTOFromEntity(userEntity)
}

func (uc *AuthUseCases) Login(ctx context.Context, email, password string) (string, *entities.AuthTokenPayload, error) {
	user, hashedPw, err := uc.userRepository.FindByEmail(ctx, email)

	if err != nil {
		return "", nil, fmt.Errorf("login failed: %w", err)
	}

	err = utils.CheckPassword(password, hashedPw)

	if user == nil || err != nil {
		return "", nil, fmt.Errorf("login failed: invalid credentials")
	}

	return uc.auth.GenerateToken(user.ID, uc.tokenDuration)

}

func (uc *AuthUseCases) ValidateToken(token string) (*entities.AuthTokenPayload, error) {
	return uc.auth.VerifyToken(token)

}

func (uc *AuthUseCases) GenerateToken(userId string, tokenDuration time.Duration) (string, *entities.AuthTokenPayload, error) {
	return uc.auth.GenerateToken(userId, tokenDuration)
}
