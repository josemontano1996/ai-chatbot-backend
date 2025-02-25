package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
)

const (
	UserDataContextKey string = "user_data"
)

func GetUserDataFromContext(ctx *gin.Context) (*entities.AuthTokenPayload, bool) {
	payload, exists := ctx.Get(UserDataContextKey)
	return payload.(*entities.AuthTokenPayload), exists
}

func GetUserIdFromContext(ctx *gin.Context) (uuid.UUID, error) {
	payload, _ := GetUserDataFromContext(ctx)
	return uuid.Parse(payload.UserId)
}

func SaveUserDataInContext(ctx *gin.Context, payload *entities.AuthTokenPayload) {
	ctx.Set(UserDataContextKey, payload)
}
