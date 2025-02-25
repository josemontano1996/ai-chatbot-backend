package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/josemontano1996/ai-chatbot-backend/internal/dto"
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/in"
	outrepo "github.com/josemontano1996/ai-chatbot-backend/internal/ports/out/repositories"
	"github.com/josemontano1996/ai-chatbot-backend/pkg/utils"
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

func (uc *UserUseCase) UpdateUser(ctx context.Context, id uuid.UUID, new_email, new_password string) (*dto.User, error) {
	hashedPw := ""
	// if the password is not null, hash it to store it in db, if not send null value
	if new_password != "" {
		hash, err := utils.HashPassword(new_password, 12)
		if err != nil {
			return nil, err
		}
		hashedPw = hash
	}

	userEntity, err := uc.userRepo.UpdateUser(ctx, id, new_email, hashedPw)
	if err != nil {
		return nil, err
	}

	return dto.NewUserDTOFromEntity(userEntity)
}
