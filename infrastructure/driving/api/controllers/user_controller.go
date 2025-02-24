package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/josemontano1996/ai-chatbot-backend/infrastructure/driven/auth"
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/in"
	outrepo "github.com/josemontano1996/ai-chatbot-backend/internal/ports/out/repositories"
)

type UserController struct {
	userUseCase in.UserUseCase
	userRepo    outrepo.UserRepository
}

func NewUserController(userUC in.UserUseCase, userRepo outrepo.UserRepository) *UserController {
	return &UserController{
		userUseCase: userUC,
		userRepo:    userRepo,
	}
}

func (c *UserController) GetUserById(ctx *gin.Context) {
	userData, exists := auth.GetUserDataFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, sendErrorPayload(errors.New("user not found")))
		return
	}

	userId, err := uuid.Parse(userData.UserId)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, sendErrorPayload(err))
		return
	}

	userDto, err := c.userUseCase.GetUserById(ctx, userId)
	fmt.Println(err, userDto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, sendErrorPayload(err))
		return
	}

	ctx.JSON(http.StatusOK, sendSuccessPayload(userDto))
}
