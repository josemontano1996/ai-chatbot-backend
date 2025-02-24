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
		authHeader := ctx.GetHeader(authorizationKey)

		// Cookies are included for authenticating only WS connections
		authCookie, _ := ctx.Cookie(authorizationKey)
		// No authentication provided, then request not accepted
		if len(authHeader) == 0 && authCookie == "" {
			err := errors.New("unauthorized")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		var accessToken string

		if authCookie != "" {
			token, err := authenticateViaCookie(authCookie)

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			accessToken = token
		} else {
			// If the authorization header is provided, authenticate via this method as it is more desirable
			token, err := authenticateViaHeader(authHeader)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}

			accessToken = token
		}

		payload, err := auth.ValidateToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

func authenticateViaHeader(authHeader string) (string, error) {
	fields := strings.Fields(authHeader)

	if len(fields) != 2 {
		err := errors.New("unauthorized")
		return "", err
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != authorizationTypeBearer {
		err := errors.New("unauthorized")
		return "", err
	}

	return fields[1], nil

}

func authenticateViaCookie(cookie string) (string, error) {
	if cookie == "" {
		return "", errors.New("not authorized")
	}

	return cookie, nil
}

func GetUserFromRequest(ctx *gin.Context) *entities.AuthTokenPayload {
	payload, _ := ctx.Get(authorizationPayloadKey)
	return payload.(*entities.AuthTokenPayload)
}
