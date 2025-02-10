package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/josemontano1996/ai-chatbot-backend/api/controller"
)

func RegisterRoutes(router *gin.Engine) (*gin.Engine, error) {

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "health checked")
	})
	router.GET("/ws", controller.PostAIController)

	return router, nil
}
