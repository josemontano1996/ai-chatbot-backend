package dto

import (
	"errors"

	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
	"github.com/josemontano1996/ai-chatbot-backend/pkg/utils"
)

type ChatHistoryDTO []ChatMessageDTO

type ChatMessageDTO struct {
	UserId  string                   `json:"userId"`
	Code    entities.ChatMessageCode `json:"code" validate:"required,oneof=1 2 3" binding:"required,oneof=1 2 3"`
	Message string                   `json:"message" binding:"required" validate:"required"`
}

func NewChatMessageDTO(code entities.ChatMessageCode, message string, userId string) (*ChatMessageDTO, error) {
	dto := ChatMessageDTO{
		Code:    code,
		Message: message,
		UserId:  userId,
	}

	err := utils.ValidateStruct(dto)

	if err != nil {
		return nil, err
	}

	return &ChatMessageDTO{
		Code:    code,
		Message: message,
	}, nil
}

func (dto *ChatMessageDTO) ToEntity() *entities.ChatMessage {
	return &entities.ChatMessage{
		Code:    dto.Code,
		Message: dto.Message,
	}
}

func ChatMessageEntityToDTO(entity *entities.ChatMessage) (*ChatMessageDTO, error) {

	if entity.Code == entities.AISystemChatMessageCode {
		return nil, errors.New("cannot convert system message to DTO")
	}

	return NewChatMessageDTO(entity.Code, entity.Message, entity.UserId)
}
