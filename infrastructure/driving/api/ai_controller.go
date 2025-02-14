package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/josemontano1996/ai-chatbot-backend/domain/entities"
	inputport "github.com/josemontano1996/ai-chatbot-backend/domain/ports/input"
	outputport "github.com/josemontano1996/ai-chatbot-backend/domain/ports/output"
	"github.com/josemontano1996/ai-chatbot-backend/infrastructure/driving/ws"
)

type AIController struct {
	aiChatUseCase         inputport.AIChatUseCase
	chatMessageRepository outputport.ChatMessageRepository
}

func NewAIController(aiChatUseCase inputport.AIChatUseCase, chatMessageRespository outputport.ChatMessageRepository) *AIController {
	return &AIController{
		aiChatUseCase:         aiChatUseCase,
		chatMessageRepository: chatMessageRespository,
	}
}

func (c *AIController) ChatWithAI(ctx *gin.Context) {
	expirationTime := 60 * time.Minute
	userID := uuid.New()
	user := &entities.User{
		ID:   userID,
		Name: "Federico",
	}

	wsClient, err := ws.NewWSClient(ctx, userID, expirationTime)

	if err != nil {
		log.Println("Error creating new WS client:", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer wsClient.Conn.Close()

	log.Println("New WS client created", time.Now())

	for {
		userMessagePayload, err := wsClient.ParseIncomingRequest()

		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				log.Println("WebSocket closed by client or server:", err) // Handle normal close scenarios
			} else {
				log.Println("Error parsing user message:", err) // Log unexpected errors
			}
			break // Exit loop on read error
		}

		userMessage, err := entities.NewUserMessage(userMessagePayload.Message)

		if err != nil {
			log.Println("Error creating new user message:", err)
			break
		}

		chatHistory, err := c.chatMessageRepository.GetChatHistory(ctx, userID.String())

		if err != nil {
			log.Println("Error getting chat history from repository:", err)
			break
		}

		chatResponse, err := c.aiChatUseCase.SendChatMessage(ctx, user, userMessage, chatHistory)

		if err != nil {
			log.Println("Error sending chat message to AI:", err)
			break
		}

		err = c.chatMessageRepository.SaveMessages(ctx, userID.String(), userMessage, chatResponse.ChatMessage)

		if err != nil {
			log.Println("Error saving messages to repository:", err)
		}

		err = wsClient.Conn.WriteJSON(chatResponse.ChatMessage)

		if err != nil {
			log.Println("Error writing response to WS client:", err)
			break
		}
	}

}
