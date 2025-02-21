package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/josemontano1996/ai-chatbot-backend/infrastructure/driving/ws"
	"github.com/josemontano1996/ai-chatbot-backend/internal/dto"
	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/in"
	outrepo "github.com/josemontano1996/ai-chatbot-backend/internal/ports/out/repositories"
)

type AIController struct {
	aiChatUseCase         in.AIChatUseCase
	chatMessageRepository outrepo.ChatMessageRepository
	ws                    ws.AIChatWSClientInterface
}

func NewAIController(aiChatUseCase in.AIChatUseCase, chatMessageRespository outrepo.ChatMessageRepository, chatWebsocket ws.AIChatWSClientInterface) *AIController {
	return &AIController{
		aiChatUseCase:         aiChatUseCase,
		chatMessageRepository: chatMessageRespository,
		ws:                    chatWebsocket,
	}
}

func (c *AIController) ChatWithAI(ctx *gin.Context) {
	// user will come from the ctx field from the middleware
	userID := "someid"
	user := &entities.User{
		ID: userID,
	}

	err := c.ws.Connect(ctx)

	if err != nil {
		log.Println("Error connecting WS:", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	defer c.ws.Disconnect()

	for {
		userMessagePayload, err := c.ws.ReadChatMessage()

		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				log.Println("WebSocket closed by client or server:", err)
				break
			}

			c.ws.SendErrorToClient(err) // Log unexpected errors
			continue
		}

		aiResponse, err := c.aiChatUseCase.SendChatMessage(ctx, user.ID, userMessagePayload.Message)

		if err != nil {
			c.ws.SendErrorToClient(err)
			continue
		}
		log.Println("AI response received:", aiResponse.ChatMessage.Message)
		chatMessageDTO, err := dto.ChatMessageEntityToDTO(aiResponse.ChatMessage)

		if err != nil {
			c.ws.SendErrorToClient(err)
		} else {
			err = c.ws.SendChatMessage(chatMessageDTO)

			if err != nil {
				c.ws.SendErrorToClient(err)
				continue
			}
		}
	}

}
