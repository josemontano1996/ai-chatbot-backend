package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

type registerRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (c *AuthController) RegisterUser(ctx *gin.Context) {
	var reqUser registerRequest

	err := ctx.ShouldBindJSON(&reqUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, sendErrorPayload(err))
		return
	}

	_, err = c.authUseCase.RegisterUser(ctx, reqUser.Email, reqUser.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, sendErrorPayload(err))
		return
	}

	ctx.JSON(http.StatusOK, sendSuccessPayload(""))
}
