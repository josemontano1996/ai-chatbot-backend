package ws

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/josemontano1996/ai-chatbot-backend/internal/dto"
	"github.com/josemontano1996/ai-chatbot-backend/pkg/utils"
)

type WSPayload[T any] struct {
	Payload T      `json:"payload" validate:"required"`
	Error   string `json:"error"`
}

type WSConfig struct {
	Ctx             *gin.Context      `binding:"omitempty"`
	ExpirationTime  time.Duration     `binding:"required"`
	ReadBufferSize  int               `binding:"required"`
	WriteBufferSize int               `binding:"required"`
	CheckOrigin     func(...any) bool `binding:"required"`
}

type WSClientInterface[T any] interface {
	ParseIncomingRequest() (*WSPayload[T], error)
	SendResposeToClient(payload *WSPayload[T]) error
	NewPayload(T, error) *WSPayload[T]
	Connect(ctx *gin.Context) error
	Disconnect() error
	SendErrorToClient(err error) error
}

type AIChatWSClientInterface interface {
	SendChatMessage(message *dto.ChatMessageDTO) error
	ReadChatMessage() (*dto.ChatMessageDTO, error)
	SendErrorToClient(err error) error
	Connect(ctx *gin.Context) error
	Disconnect() error
}

func NewWSConfig(readBufferSize int, writeBufferSize int, expiration time.Duration) (*WSConfig, error) {
	config := &WSConfig{
		ExpirationTime:  expiration,
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
		CheckOrigin: func(...any) bool {
			return true
		},
	}

	err := utils.ValidateStruct(config)

	return config, err
}
