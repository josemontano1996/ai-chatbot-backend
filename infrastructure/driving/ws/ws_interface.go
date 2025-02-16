package ws

import (
	"time"

	"github.com/gin-gonic/gin"
)

type WSPayload[T any] struct {
	Payload T `json:"payload" validate:"required"`
}

type WSConfig struct {
	Ctx             *gin.Context
	ExpirationTime  time.Duration
	ReadBufferSize  int
	WriteBufferSize int
	CheckOrigin     func(...any) bool
}

type WSClientInterface[T any] interface {
	ParseIncomingRequest() (*WSPayload[T], error)
	SendResposeToClient(payload *WSPayload[T]) error
	NewPayload(T) *WSPayload[T]
	Connect(config WSConfig) error
	Disconnect() error
}
