package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/in"
)

const (
	authorizationKey        = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(auth in.AuthUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authozarionHeader := ctx.GetHeader(authorizationKey)
		if len(authozarionHeader) == 0 {
			err := errors.New("unauthorized")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		fields := strings.Fields(authozarionHeader)
		if len(fields) < 2 {
			err := errors.New("unauthorized")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := errors.New("unauthorized")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		accessToken := fields[1]
		payload, err := auth.ValidateToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

func GetUserFromRequest(ctx *gin.Context) *entities.AuthTokenPayload {
	payload, _ := ctx.Get(authorizationPayloadKey)
	return payload.(*entities.AuthTokenPayload)
}
