package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	chatws "github.com/josemontano1996/ai-chatbot-backend/infrastructure/driving/ws/chat"
	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/in"
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/out"
)

type AIController struct {
	aiChatUseCase         in.AIChatUseCase
	chatMessageRepository out.ChatMessageRepository
	ws                    chatws.AIChatWSClientInterface
}

func NewAIController(aiChatUseCase in.AIChatUseCase, chatMessageRespository out.ChatMessageRepository, chatWebsocket chatws.AIChatWSClientInterface) *AIController {
	return &AIController{
		aiChatUseCase:         aiChatUseCase,
		chatMessageRepository: chatMessageRespository,
		ws:                    chatWebsocket,
	}
}

func (c *AIController) ChatWithAI(ctx *gin.Context) {
	expirationTime := 60 * time.Minute
	userID := uuid.New()

	// user will come from the ctx field from the middleware
	user := &entities.User{
		ID:   userID,
		Name: "Federico",
	}

	wsClient, err := c.ws.NewWSClient(ctx, userID, expirationTime)

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
