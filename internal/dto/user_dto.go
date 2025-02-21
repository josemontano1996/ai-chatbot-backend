package dto

import (
	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
	"github.com/josemontano1996/ai-chatbot-backend/pkg/utils"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email" validate:"required,email" bindig:"required,email"`
	Password string `json:"password"`
}


func NewUserDTOFromEntity(user *entities.User) (*User, error) {
	userDTO := &User{
		ID:    user.ID,
		Email: user.Email,
	}
	err := utils.ValidateStruct(userDTO)

	if err != nil {
		return nil, err
	}

	return userDTO, nil
}
