package chatws

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/josemontano1996/ai-chatbot-backend/infrastructure/driving/ws"
	"github.com/josemontano1996/ai-chatbot-backend/internal/dto"
)

type AIChatWSClient struct {
	client ws.WSClientInterface[dto.ChatMessageDTO]
}

func NewAIChatWSClient(config ws.WSConfig) (ws.AIChatWSClientInterface, error) {

	client, err := ws.NewGorillaWSClient[dto.ChatMessageDTO](config)

	if err != nil {
		return nil, err
	}

	return &AIChatWSClient{
		client: client,
	}, nil
}

func (c *AIChatWSClient) ReadChatMessage() (*dto.ChatMessageDTO, error) {
	payloadWrapper, err := c.client.ParseIncomingRequest()

	if err != nil {
		return nil, err
	}

	chatMessagePayload := payloadWrapper.Payload

	return &chatMessagePayload, nil
}

func (c *AIChatWSClient) SendChatMessage(message *dto.ChatMessageDTO) error {
	if message == nil {
		return errors.New("nil pointer for message")
	}
	payload := c.client.NewPayload(*message, nil)
	return c.client.SendResposeToClient(payload)
}


func (c *AIChatWSClient) Connect(ctx *gin.Context) error {
	return c.client.Connect(ctx)
}
func (c *AIChatWSClient) Disconnect() error {
	return c.client.Disconnect()
}
func (c *AIChatWSClient) SendErrorToClient(err error) error {
	return c.client.SendErrorToClient(err)
}
