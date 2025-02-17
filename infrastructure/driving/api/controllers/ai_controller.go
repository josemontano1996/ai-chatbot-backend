package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/josemontano1996/ai-chatbot-backend/infrastructure/driving/ws"
	"github.com/josemontano1996/ai-chatbot-backend/internal/dto"
	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/in"
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/out"
)

type AIController struct {
	aiChatUseCase         in.AIChatUseCase
	chatMessageRepository out.ChatMessageRepository
	ws                    ws.AIChatWSClientInterface
}

func NewAIController(aiChatUseCase in.AIChatUseCase, chatMessageRespository out.ChatMessageRepository, chatWebsocket ws.AIChatWSClientInterface) *AIController {
	return &AIController{
		aiChatUseCase:         aiChatUseCase,
		chatMessageRepository: chatMessageRespository,
		ws:                    chatWebsocket,
	}
}

func (c *AIController) ChatWithAI(ctx *gin.Context) {
	// user will come from the ctx field from the middleware
	userID := uuid.New()
	user := &entities.User{
		ID:   userID,
		Name: "Federico",
	}

	err := c.ws.Connect(ctx)

	if err != nil {
		log.Println("Error connecting WS:", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	defer c.ws.Disconnect()

	log.Println("New WS client created", time.Now())

	for {
		userMessagePayload, err := c.ws.ReadChatMessage()

		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				log.Println("WebSocket closed by client or server:", err) // Handle normal close scenarios
			} else {
				log.Println("Error parsing user message:", err) // Log unexpected errors
			}
			break // Exit loop on read error
		}

		aiResponse, err := c.aiChatUseCase.SendChatMessage(ctx, user.ID.String(), userMessagePayload.Message)

		if err != nil {
			log.Println("Error sending message to AI:", err)
			break
		}

		chatMessageDTO, err := dto.ChatMessageEntityToDTO(aiResponse.ChatMessage)

		if err != nil {
			log.Println("Error converting entity to DTO:", err)
		} else {
			err = c.ws.SendChatMessage(chatMessageDTO)

			if err != nil {
				log.Println("Error writing response to WS client:", err)
				break
			}
		}
	}

}
