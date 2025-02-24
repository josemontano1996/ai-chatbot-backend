package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/josemontano1996/ai-chatbot-backend/internal/dto"
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/in"
	outrepo "github.com/josemontano1996/ai-chatbot-backend/internal/ports/out/repositories"
)

type UserUseCase struct {
	userRepo outrepo.UserRepository
}

func NewUserUseCase(userRepo outrepo.UserRepository) in.UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (uc *UserUseCase) GetUserById(ctx context.Context, id uuid.UUID) (*dto.User, error) {
	userEntity, _, err := uc.userRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	userDto, err := dto.NewUserDTOFromEntity(userEntity)
	if err != nil {
		return nil, err
	}

	return userDto, err
}
