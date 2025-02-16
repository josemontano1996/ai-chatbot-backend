package usecases

import (
	"context"

	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/in"
	"github.com/josemontano1996/ai-chatbot-backend/internal/ports/out"
)

type AIChatUseCase[T any] struct {
	aiProvider  out.AIProvider[T]
	messageRepo out.ChatMessageRepository
}

func NewAIChatUseCase[T any](aiProvider out.AIProvider[T], messageRepo out.ChatMessageRepository) in.AIChatUseCase {
	return &AIChatUseCase[T]{
		aiProvider:  aiProvider,
		messageRepo: messageRepo,
	}
}

func (uc *AIChatUseCase[T]) SendChatMessage(ctx context.Context, user *entities.User, userMessage *entities.ChatMessage, history *entities.ChatHistory) (*in.AIChatResponse, error) {

	aiResponse, _, err := uc.aiProvider.SendMessage(ctx, userMessage, history)

	if err != nil {
		return &in.AIChatResponse{}, err
	}

	// TODO: Implement logic to substract tokens from user's account
	// metadata.TokensSpent
	err = uc.messageRepo.SaveMessages(ctx, user.ID.String(), userMessage, aiResponse.ChatMessage)

	if err != nil {
		return &in.AIChatResponse{}, err
	}

	return aiResponse, nil
}
