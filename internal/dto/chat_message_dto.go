package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
)

type ChatHistoryDTO []ChatMessageDTO

type ChatMessageDTO struct {
	Code    entities.ChatMessageCode `json:"code" validate:"required,oneof=1 2 3" binding:"required,oneof=1 2 3"`
	Message string                   `json:"message" binding:"required" validate:"required"`
}

func NewChatMessageDTO(code entities.ChatMessageCode, message string) (*ChatMessageDTO, error) {
	dto := ChatMessageDTO{
		Code:    code,
		Message: message,
	}

	v := validator.New()

	err := v.Struct(dto)

	if err != nil {
		return nil, err
	}

	return &ChatMessageDTO{
		Code:    code,
		Message: message,
	}, nil
}

func NewUserChatMessageDTO(message string) (*ChatMessageDTO, error) {
	return NewChatMessageDTO(entities.UserChatMessageCode, message)
}

func (dto *ChatMessageDTO) ToEntity() *entities.ChatMessage {
	return &entities.ChatMessage{
		Code:    dto.Code,
		Message: dto.Message,
	}
}

func ChatMessageEntityToDTO(entity *entities.ChatMessage) (*ChatMessageDTO, error) {

	if entity.Code == entities.AISystemChatMessageCode {
		return &ChatMessageDTO{}, nil
	}

	return NewChatMessageDTO(entity.Code, entity.Message)
}
