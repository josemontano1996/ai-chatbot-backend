package controller

import (
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/in"
	outrepo "github.com/josemontano1996/ai-chatbot-backend/internal/ports/out/repositories"
)

type AuthController struct {
	authUseCase in.AuthUseCase
	userRepo    outrepo.UserRepository
}

func NewAuthController(authUseCase in.AuthUseCase, userRepo outrepo.UserRepository) *AuthController {
	return &AuthController{
		authUseCase: authUseCase,
		userRepo:    userRepo,
	}
}
