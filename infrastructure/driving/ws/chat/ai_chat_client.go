package chatws

import (
	"errors"

	"github.com/josemontano1996/ai-chatbot-backend/infrastructure/driving/ws"
	"github.com/josemontano1996/ai-chatbot-backend/internal/dto"
)




type AIChatWSClient struct {
	client ws.WSClientInterface[dto.ChatMessageDTO]
}

func NewAIChatWSClient() ws.AIChatWSClientInterface {
	return &AIChatWSClient{
		client: ws.NewGorillaWSClient[dto.ChatMessageDTO](),
	}
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
	payload := c.client.NewPayload(*message)
	return c.client.SendResposeToClient(payload)
}

func (c *AIChatWSClient) Connect(config ws.WSConfig) error {
	return c.client.Connect(config)
}
func (c *AIChatWSClient) Disconnect() error {
	return c.client.Disconnect()
}
